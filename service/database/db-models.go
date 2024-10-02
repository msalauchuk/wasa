package database

import "time"

/*
// Photo structure for the database
type Photo struct {
	Comments int       // Number of comments of a photo
	Likes    int       // Number of likes of a photo
	Owner    string    // Unique id of the owner
	PhotoId  int       // Unique id of the photo
	Date     time.Time // Date in which the photo was uploaded
}
*/

// User defines a user structure.
type User struct {
	IdUser string `json:"user_id"` // User's unique ID
}

// CompleteUser includes detailed user information.
type CompleteUser struct {
	IdUser   string `json:"user_id"`  // User's unique ID
	Nickname string `json:"nickname"` // Nickname of the user
}

// Photo holds details about a photo in the database.
type Photo struct {
	PhotoId  int               `json:"photo_id"` // Unique ID of the photo
	Owner    string            `json:"owner"`    // Unique ID of the owner
	Date     time.Time         `json:"date"`     // Date the photo was uploaded
	Comments []CompleteComment `json:"comments"` // Array of comments on the photo
	Likes    []CompleteUser    `json:"likes"`    // Array of users who liked the photo
}

// CompleteComment includes all details for a comment.
type CompleteComment struct {
	IdComment int64  `json:"comment_id"` // Identifier of the comment
	IdPhoto   int64  `json:"photo_id"`   // Photo unique ID
	IdUser    string `json:"user_id"`    // User's unique ID
	Nickname  string `json:"nickname"`   // Nickname of the user
	Comment   string `json:"comment"`    // Content of the comment
}

// PhotoId structure for the database
type PhotoId struct {
	IdPhoto int64 `json:"photo_id"` // Photo unique id
}

// Nickname structure for the database
type Nickname struct {
	Nickname string `json:"nickname"` // Nickname of a user
}

// Comment structure for the database
type Comment struct {
	Comment string `json:"comment"` // Comment content
}

// CommentId structure for the database
type CommentId struct {
	IdComment int64 `json:"comment_id"` // Identifier of a comment
}
