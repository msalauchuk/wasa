package database

// GetPhotosList retrieves a list of photos for a specific user, ordered by date in descending order.
// The photos are only returned if the requesting user is not banned by the target user.
// Returns:
// - A slice of Photo structs containing the photos' metadata, comments, and likes.
// - An error if something goes wrong during database queries.
func (db *appdbimpl) GetPhotosList(requestingUser User, targetUser User) ([]Photo, error) {

	rows, err := db.c.Query(""+
		"SELECT * FROM photos "+
		"WHERE id_user = ? "+
		"ORDER BY date DESC", targetUser.IdUser)
	if err != nil {
		return nil, err
	}

	// Wait for the function to finish before closing rows
	defer func() { _ = rows.Close() }()

	// Initialize slice to hold photos
	var photos []Photo

	// Iterate over the result set
	for rows.Next() {
		var photo Photo
		err = rows.Scan(&photo.PhotoId, &photo.Owner, &photo.Date)
		if err != nil {
			return nil, err
		}

		// Fetch the list of likes for this photo
		likes, err := db.GetLikesList(requestingUser, targetUser, PhotoId{IdPhoto: int64(photo.PhotoId)}) // Old: GetLikesLen
		if err != nil {
			return nil, err
		}
		photo.Likes = likes

		// Fetch the list of comments for this photo
		comments, err := db.GetCompleteCommentsList(requestingUser, targetUser, PhotoId{IdPhoto: int64(photo.PhotoId)}) // Old: GetCommentsLen
		if err != nil {
			return nil, err
		}
		photo.Comments = comments

		// Add photo to the result slice
		photos = append(photos, photo)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return photos, nil
}

// GetPhoto retrieves a specific photo's metadata from the database if the requesting user is not banned by the photo's owner.
// Returns:
// - A Photo struct containing the photo's metadata.
// - An error if the user is banned or if another issue occurs.
func (db *appdbimpl) GetPhoto(requestingUser User, targetPhoto PhotoId) (Photo, error) {

	var photo Photo
	err := db.c.QueryRow(""+
		"SELECT * FROM photos"+
		"WHERE (id_photo = ?) "+
		"AND id_user NOT IN (SELECT banner FROM banned_user WHERE banned = ?)",
		targetPhoto.IdPhoto, requestingUser.IdUser).Scan(&photo)

	if err != nil {
		return Photo{}, ErrUserBanned
	}

	return photo, nil

}

// CreatePhoto inserts a new photo into the database and returns its unique photo ID.
// Parameters:
// - p: the Photo struct containing the owner ID and date of the photo.
// Returns:
// - The unique ID of the inserted photo.
// - An error if the insertion fails.
func (db *appdbimpl) CreatePhoto(p Photo) (int64, error) {

	res, err := db.c.Exec("INSERT INTO photos (id_user,date) VALUES (?,?)",
		p.Owner, p.Date)

	if err != nil {
		// Error executing query
		return -1, err
	}

	photoId, err := res.LastInsertId()
	if err != nil {
		// Error getting id returned by last db operation (photoId)
		return -1, err
	}

	return photoId, nil
}

// RemovePhoto deletes a photo from the database, ensuring only the owner can delete it.
// Parameters:
// - owner: the user who owns the photo.
// - p: the ID of the photo to be deleted.
// Returns:
// - An error if the deletion fails.
func (db *appdbimpl) RemovePhoto(owner User, p PhotoId) error {

	_, err := db.c.Exec(""+
		"DELETE FROM photos WHERE id_user = ? "+
		"AND id_photo = ? ",
		owner.IdUser, p.IdPhoto)
	if err != nil {
		// Error during the execution of the query
		return err
	}

	return nil
}

// CheckPhotoExistence checks if a specific photo exists in the database.
//
// Parameters:
// - targetPhoto: the ID of the photo being checked.
//
// Returns:
// - true if the photo exists, false otherwise.
// - An error if there is a failure during the query.
func (db *appdbimpl) CheckPhotoExistence(targetPhoto PhotoId) (bool, error) {

	var rows int
	err := db.c.QueryRow(""+
		"SELECT COUNT(*) FROM photos"+
		"WHERE (id_photo = ?)", targetPhoto.IdPhoto).Scan(&rows)
	if err != nil {
		return false, err
	}

	if rows == 0 {
		return false, nil
	}
	return true, nil

}
