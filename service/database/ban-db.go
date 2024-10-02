package database

// BanUser bans the user `banned` by adding an entry to the `banned_users` table.
// `banner` is the user who initiates the ban.
// Returns an error if the operation fails.
func (db *appdbimpl) BanUser(banner User, banned User) error {

	_, err := db.c.Exec("INSERT INTO banned_users (banner,banned) VALUES (?, ?)", banner.IdUser, banned.IdUser)
	if err != nil {
		return err
	}

	return nil
}

// UnbanUser removes the user `banned` from the banned list of `banner` by deleting the entry from the `banned_users` table.
// Returns an error if the operation fails.
func (db *appdbimpl) UnbanUser(banner User, banned User) error {

	_, err := db.c.Exec("DELETE FROM banned_users WHERE (banner = ? AND banned = ?)", banner.IdUser, banned.IdUser)
	if err != nil {
		return err
	}

	return nil
}

// BannedUserCheck checks if `requestingUser` has been banned by `targetUser`.
// Returns true if the user is banned, false otherwise.
// Also returns an error if the query execution fails.
func (db *appdbimpl) BannedUserCheck(requestingUser User, targetUser User) (bool, error) {

	var cnt int
	err := db.c.QueryRow("SELECT COUNT(*) FROM banned_users WHERE banned = ? AND banner = ?",
		requestingUser.IdUser, targetUser.IdUser).Scan(&cnt)

	if err != nil {
		// COUNT(*) always returns a row, so this error should be rare.
		return true, err
	}

	// If the count is 1, it means the user is banned.
	if cnt == 1 {
		return true, nil
	}
	return false, nil
}
