package api

import (
	"time"
	"wasaphoto-1764024/service/database"
)

const ImgFormatErrorMsg = "images must be jpeg or png"

// JSONErrorMsg JSON Error Structure
type JSONErrorMsg struct {
	Message string `json:"message"` // Error messages
}

// Photo represents a photo in the API layer
type Photo struct {
	Comments []database.CompleteComment `json:"comments"` // Number of comments of a photo
	Likes    []database.CompleteUser    `json:"likes"`    // Number of likes of a photo
	Owner    string                     `json:"owner"`    // Unique id of the owner
	PhotoId  int                        `json:"photo_id"` // Unique id of the photo
	Date     time.Time                  `json:"date"`     // Date in which the photo was uploaded
}

// User represents a user in the API layer
type User struct {
	IdUser string `json:"user_id"` // User's unique id
}

// PhotoId structure for the APIs
type PhotoId struct {
	IdPhoto int64 `json:"photo_id"` // Photo unique id
}

// Nickname structure for the APIs
type Nickname struct {
	Nickname string `json:"nickname"` // Nickname of a user
}

// Comment structure for the APIs
type Comment struct {
	Comment string `json:"comment"` // Comment content
}

// CommentId structure for the APIs
type CommentId struct {
	IdComment int64 `json:"comment_id"` // Identifier of a comment
}

// CompleteComment structure for the APIs
type CompleteComment struct {
	IdComment int64  `json:"comment_id"` // Identifier of a comment
	IdPhoto   int64  `json:"photo_id"`   // Photo unique id
	IdUser    string `json:"user_id"`    // User's unique id
	Nickname  string `json:"nickname"`   // Nickname of a user
	Comment   string `json:"comment"`    // Comment content
}

// CompleteProfile structure for the APIs
type CompleteProfile struct {
	Name      string           `json:"user_id"`
	Nickname  string           `json:"nickname"`
	Followers []database.User  `json:"followers"`
	Following []database.User  `json:"following"`
	Posts     []database.Photo `json:"posts"`
}

// ToDatabase Maps the User from the api package to a User of the database package
func (u User) ToDatabase() database.User {
	return database.User{
		IdUser: u.IdUser,
	}
}

// ToDatabase Maps the Photo from the api package to a Photo of the database package
func (p Photo) ToDatabase() database.Photo {
	return database.Photo{
		Comments: p.Comments,
		Likes:    p.Likes,
		Owner:    p.Owner,
		PhotoId:  p.PhotoId,
		Date:     p.Date,
	}
}

// ToDatabase Maps the PhotoId from the api package to a PhotoId of the database package
func (p PhotoId) ToDatabase() database.PhotoId {
	return database.PhotoId{
		IdPhoto: p.IdPhoto,
	}
}

// ToDatabase Maps the Nickname from the api package to a Nickname of the database package
func (n Nickname) ToDatabase() database.Nickname {
	return database.Nickname{
		Nickname: n.Nickname,
	}
}

// ToDatabase Maps the Comment from the api package to a Comment of the database package
func (c Comment) ToDatabase() database.Comment {
	return database.Comment{
		Comment: c.Comment,
	}
}

// ToDatabase Maps the CommentId from the api package to a CommentId of the database package
func (c CommentId) ToDatabase() database.CommentId {
	return database.CommentId{
		IdComment: c.IdComment,
	}
}

// ToDatabase Maps the CompleteComment from the api package to a CompleteComment of the database package
func (cc CompleteComment) ToDatabase() database.CompleteComment {
	return database.CompleteComment{
		IdComment: cc.IdComment,
		IdPhoto:   cc.IdPhoto,
		IdUser:    cc.IdUser,
		Comment:   cc.Comment,
	}
}
