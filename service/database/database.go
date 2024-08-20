/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.
To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.
For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// Custom errors for database operations
var (
	ErrPhotoDoesNotExist = errors.New("photo does not exist")
	ErrUserBanned        = errors.New("user is banned")
)

// Constants for operation limits
const PhotosPerUserHome = 3

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

	// User management
	CreateUser(user User) error
	ChangeNickname(user User, nickname Nickname) error
	SearchUser(searcher User, userToSearch User) ([]CompleteUser, error)

	// Photo management
	CreatePhoto(Photo) (int64, error)
	LikePhoto(PhotoId, User) error
	UnlikePhoto(PhotoId, User) error
	CommentPhoto(PhotoId, User, Comment) (int64, error)
	UncommentPhoto(PhotoId, User, CommentId) error
	RemovePhoto(User, PhotoId) error

	// Relationship management
	FollowUser(a User, b User) error
	UnfollowUser(a User, b User) error
	BanUser(a User, b User) error
	UnbanUser(a User, b User) error

	// Social and streaming
	GetStream(User) ([]Photo, error)
	GetFollowers(User) ([]User, error)
	GetFollowing(User) ([]User, error)
	GetPhotosList(a User, b User) ([]Photo, error)
	UncommentPhotoAuthor(PhotoId, CommentId) error

	// Utility methods
	GetNickname(User) (string, error)
	BannedUserCheck(a User, b User) (bool, error)
	CheckUser(a User) (bool, error)
	CheckPhotoExistence(p PhotoId) (bool, error)

	// Ping checks whether the database is available or not (in that case, an error will be returned)
	Ping() error
	// Gets the user by ID. Returns the user and an error.
	GetUserByID(id string) (*User, error)
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if the database is live
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database not accessible: %w", err)
	}

	requiredTables := []string{"users", "photos", "likes", "comments", "banned_users", "followers"}
	if err := ensureDatabaseTables(db, requiredTables); err != nil {
		return nil, err
	}

	return &appdbimpl{
		c: db,
	}, nil
}

// Check the existence of tables
func ensureDatabaseTables(db *sql.DB, requiredTables []string) error {
	var name string
	for _, table := range requiredTables {
		query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?;`
		err := db.QueryRow(query, table).Scan(&name)
		if errors.Is(err, sql.ErrNoRows) {
			// If any table does not exist, create all necessary tables
			if err := createDatabase(db); err != nil {
				return fmt.Errorf("error creating database structure for %s: %w", table, err)
			}
			break // Exit after tables are created
		} else if err != nil {
			return fmt.Errorf("error checking table %s: %w", table, err)
		}
	}
	return nil
}
func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

// Define table creation SQL statements
func createDatabase(db *sql.DB) error {
	schemaStatements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id_user VARCHAR(16) NOT NULL PRIMARY KEY,
			nickname VARCHAR(16) NOT NULL
			);`,
		`CREATE TABLE IF NOT EXISTS photos (
			id_photo INTEGER PRIMARY KEY AUTOINCREMENT,
			id_user VARCHAR(16) NOT NULL,
			date DATETIME NOT NULL,
			FOREIGN KEY(id_user) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS  likes (
			id_photo INTEGER NOT NULL,
			id_user VARCHAR(16) NOT NULL,
			PRIMARY KEY (id_photo,id_user),
			FOREIGN KEY(id_photo) REFERENCES photos (id_photo) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS comments (
			id_comment INTEGER PRIMARY KEY AUTOINCREMENT,
			id_photo INTEGER NOT NULL,
			id_user VARCHAR(16) NOT NULL,
			comment VARCHAR(30) NOT NULL,
			FOREIGN KEY(id_photo) REFERENCES photos (id_photo) ON DELETE CASCADE,
			FOREIGN KEY(id_user) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS banned_users (
			banner VARCHAR(16) NOT NULL,
			banned VARCHAR(16) NOT NULL,
			PRIMARY KEY (banner,banned),
			FOREIGN KEY(banner) REFERENCES users (id_user) ON DELETE CASCADE,
			FOREIGN KEY(banned) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
		`CREATE TABLE IF NOT EXISTS followers(
			follower VARCHAR(16) NOT NULL,
			followed VARCHAR(16) NOT NULL,
			PRIMARY KEY (follower,followed),
			FOREIGN KEY(follower) REFERENCES users (id_user) ON DELETE CASCADE,
			FOREIGN KEY(followed) REFERENCES users (id_user) ON DELETE CASCADE
			);`,
	}

	// Execute each schema statement within the transaction
	for _, sqlStmt := range schemaStatements {
		if _, err := db.Exec(sqlStmt); err != nil {
			return fmt.Errorf("error creating table: %w", err)
		}
	}
	return nil
}
