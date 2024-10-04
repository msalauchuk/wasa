package database

// GetCompleteCommentsList retrieves the list of comments for a photo, excluding comments from users who have banned the requesting user.
func (db *appdbimpl) GetCompleteCommentsList(requestingUser User, requestedUser User, photo PhotoId) ([]CompleteComment, error) {

	// Execute SQL query to retrieve comments while excluding banned users.
	rows, err := db.c.Query("SELECT * FROM comments WHERE id_photo = ? "+
		"AND id_user NOT IN (SELECT banned FROM banned_users WHERE banner = ? OR banner = ?) "+
		"AND id_user NOT IN (SELECT banner FROM banned_users WHERE banned = ?)",
		photo.IdPhoto, requestingUser.IdUser, requestedUser.IdUser, requestingUser.IdUser)
	if err != nil {
		return nil, err
	}

	// Wait for the function to finish before closing rows
	defer func() { _ = rows.Close() }()

	// Parse result set into comments list.
	var comments []CompleteComment
	for rows.Next() {
		var comment CompleteComment
		if err = rows.Scan(&comment.IdComment, &comment.IdPhoto, &comment.IdUser, &comment.Comment); err != nil {
			return nil, err
		}

		// Retrieve the nickname of the user who made the comment.
		nickname, err := db.GetNickname(User{IdUser: comment.IdUser})
		if err != nil {
			return nil, err
		}
		comment.Nickname = nickname
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

// CommentPhoto adds a comment from a user to a photo and returns the new comment's ID.

func (db *appdbimpl) CommentPhoto(p PhotoId, u User, c Comment) (int64, error) {

	res, err := db.c.Exec("INSERT INTO comments (id_photo,id_user,comment) VALUES (?, ?, ?)",
		p.IdPhoto, u.IdUser, c.Comment)
	if err != nil {
		// Error executing query
		return -1, err
	}

	commentId, err := res.LastInsertId()
	if err != nil {
		// Error getting id returned by last db operation (commentId)
		return -1, err
	}

	return commentId, nil
}

// UncommentPhoto removes a user's comment from a photo if they are the author.

func (db *appdbimpl) UncommentPhoto(p PhotoId, u User, c CommentId) error {

	_, err := db.c.Exec("DELETE FROM comments WHERE (id_photo = ? AND id_user = ? AND id_comment = ?)",
		p.IdPhoto, u.IdUser, c.IdComment)
	if err != nil {
		return err
	}

	return nil
}

// UncommentPhotoByAuthor allows the photo's author to delete any comment on their photo.

func (db *appdbimpl) UncommentPhotoAuthor(p PhotoId, c CommentId) error {

	_, err := db.c.Exec("DELETE FROM comments WHERE (id_photo = ? AND id_comment = ?)",
		p.IdPhoto, c.IdComment)
	if err != nil {
		return err
	}

	return nil
}
