package database

// GetLikesList retrieves the list of users who liked a given photo.
// The users are filtered:
// - No users who have banned the requesting or requested user.
// - No users banned by the requesting user.
// Returns:
// - A list of CompleteUser objects representing the users who liked the photo.
// - An error if something goes wrong.
func (db *appdbimpl) GetLikesList(requestingUser User, requestedUser User, photo PhotoId) ([]CompleteUser, error) {

	rows, err := db.c.Query("SELECT id_user FROM likes WHERE id_photo = ? "+
		"AND id_user NOT IN (SELECT banned FROM banned_users WHERE banner = ? OR banner = ?) "+
		"AND id_user NOT IN (SELECT banner FROM banned_users WHERE banned = ?)",
		photo.IdPhoto, requestingUser.IdUser, requestedUser.IdUser, requestingUser.IdUser)
	if err != nil {
		return nil, err
	}
	// Wait for the function to finish before closing rows.
	defer func() { _ = rows.Close() }()

	// Fetches a list of users who liked a photo and includes their nicknames
	var likes []CompleteUser
	for rows.Next() {
		var user CompleteUser
		err = rows.Scan(&user.IdUser)
		if err != nil {
			return nil, err
		}

		// retrieve the nickname for the user with the extracted IdUser.
		nickname, err := db.GetNickname(User{IdUser: user.IdUser})
		if err != nil {
			return nil, err
		}
		user.Nickname = nickname

		likes = append(likes, user)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return likes, nil
}

// UnlikePhoto removes a like from a user on a photo.
func (db *appdbimpl) UnlikePhoto(p PhotoId, u User) error {

	_, err := db.c.Exec("DELETE FROM likes WHERE(id_photo = ? AND id_user = ?)", p.IdPhoto, u.IdUser)
	if err != nil {
		return err
	}

	return nil
}

// LikePhoto adds a like from a user to a photo.

func (db *appdbimpl) LikePhoto(p PhotoId, u User) error {

	_, err := db.c.Exec("INSERT INTO likes (id_photo,id_user) VALUES (?, ?)", p.IdPhoto, u.IdUser)
	if err != nil {
		return err
	}

	return nil
}
