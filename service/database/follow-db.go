package database

// GetFollowers retrieves the list of followers for a given user.
// It returns a slice of User and an error, if any occurs.
func (db *appdbimpl) GetFollowers(requestingUser User) ([]User, error) {

	// Execute the query to find all followers of the requested user.
	rows, err := db.c.Query("SELECT follower FROM followers WHERE followed = ?", requestingUser.IdUser)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }() // Ensure rows are closed properly.

	var followers []User

	// Iterate over the result set and append each follower to the slice.
	for rows.Next() {
		var follower User
		if err = rows.Scan(&follower.IdUser); err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}

	// Check if there was any error during the iteration.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return followers, nil
}

// GetFollowing retrieves the list of users that the given user is following.
// It returns a slice of User and an error, if any occurs.
func (db *appdbimpl) GetFollowing(requestinUser User) ([]User, error) {

	rows, err := db.c.Query("SELECT followed FROM followers WHERE follower = ?", requestinUser.IdUser)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	// Iterate over the result set and append each following to the slice.
	var following []User
	for rows.Next() {
		var folowed User
		err = rows.Scan(&folowed.IdUser)
		if err != nil {
			return nil, err
		}
		following = append(following, folowed)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return following, nil
}

// FollowUser adds a new relationship where `follower` starts following `followed`.
// It returns an error if the operation fails.
func (db *appdbimpl) FollowUser(follower User, followed User) error {

	_, err := db.c.Exec("INSERT INTO followers (follower,followed) VALUES (?, ?)",
		follower.IdUser, followed.IdUser)
	if err != nil {
		return err
	}

	return nil
}

// UnfollowUser removes the relationship where `follower` is following `followed`.
// It returns an error if the operation fails.
func (db *appdbimpl) UnfollowUser(follower User, followed User) error {

	_, err := db.c.Exec("DELETE FROM followers WHERE(follower = ? AND followed = ?)",
		follower.IdUser, followed.IdUser)
	if err != nil {
		return err
	}

	return nil
}
