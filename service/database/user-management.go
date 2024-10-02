package database

import (
	"database/sql"
)

// GetStream retrieves the most recent photos posted by users followed by the specified user.
// The photos are returned in descending order by date.

func (db *appdbimpl) GetStream(user User) ([]Photo, error) {

	rows, err := db.c.Query(`SELECT * FROM photos WHERE id_user IN (SELECT followed FROM followers WHERE follower = ?) ORDER BY date DESC`,
		user.IdUser)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }() // Ensure rows are closed once function completes

	// Read all the users in the resulset
	var res []Photo
	for rows.Next() {
		var photo Photo
		// Scan the photo fields from the current row into the photo struct
		err = rows.Scan(&photo.PhotoId, &photo.Owner, &photo.Date) //  &photo.Comments, &photo.Likes,
		if err != nil {
			return nil, err
		}

		// Add the photo to the results
		res = append(res, photo)
	}
	// Check for errors encountered during row iteration
	if rows.Err() != nil {
		return nil, err
	}

	return res, nil
}

// CreateUser inserts a new user record into the database with the provided user data.
func (db *appdbimpl) CreateUser(u User) error {

	_, err := db.c.Exec("INSERT INTO users (id_user,nickname) VALUES (?, ?)",
		u.IdUser, u.IdUser)

	if err != nil {
		return err
	}

	return nil
}

// GetUserByID retrieves a user by their unique identifier (id_user).
// If no user is found, it returns nil without an error.
func (db *appdbimpl) GetUserByID(id string) (*User, error) {
	var u User
	err := db.c.QueryRow("SELECT id_user FROM users WHERE id_user = ?", id).Scan(&u.IdUser)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil, nil if user is not found
		} // User not found
		return nil, err
	}
	return &u, nil
}

// CheckUser verifies if a user with the specified id_user exists in the database.
func (db *appdbimpl) CheckUser(targetUser User) (bool, error) {

	var cnt int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users WHERE id_user = ?",
		targetUser.IdUser).Scan(&cnt)

	if err != nil {
		// Count always returns a row thanks to COUNT(*), so this situation should not happen
		return true, err
	}

	if cnt == 1 {
		return true, nil
	} // User exists if count is 1
	return false, nil
}

// SearchUser performs a search for users by matching the user's identifier or nickname.
// It excludes users who have been banned by the searcher and returns a list of matching users.
func (db *appdbimpl) SearchUser(searcher User, userToSearch User) ([]CompleteUser, error) {

	rows, err := db.c.Query("SELECT * FROM users WHERE ((id_user LIKE ?) OR (nickname LIKE ?)) AND id_user NOT IN (SELECT banner FROM banned_users WHERE banned = ?)",
		userToSearch.IdUser+"%", userToSearch.IdUser+"%", searcher.IdUser)
	if err != nil {
		return nil, err
	}
	// Wait for the function to finish before closing rows.
	defer func() { _ = rows.Close() }() // Ensure rows are closed after function finishes

	// Read all the users in the resulset.
	var res []CompleteUser
	for rows.Next() {
		var user CompleteUser
		err = rows.Scan(&user.IdUser, &user.Nickname)
		if err != nil {
			return nil, err
		}

		// Add the user to the results
		res = append(res, user)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return res, nil
}

// GetNickname retrieves the nickname of the specified user by their unique identifier (id_user).
func (db *appdbimpl) GetNickname(user User) (string, error) {

	var nickname string

	err := db.c.QueryRow(`SELECT nickname FROM users WHERE id_user = ?`, user.IdUser).Scan(&nickname)
	if err != nil {
		// Error during the execution of the query
		return "", err
	}
	return nickname, nil
}

// ChangeNickname updates the nickname of the specified user.
func (db *appdbimpl) ChangeNickname(user User, newNickname Nickname) error {

	_, err := db.c.Exec(`UPDATE users SET nickname = ? WHERE id_user = ?`, newNickname.Nickname, user.IdUser)
	if err != nil {
		// Error during the execution of the query
		return err
	}
	return nil
}
