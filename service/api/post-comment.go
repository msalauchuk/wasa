package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wasaphoto-1764024/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// postComment adds a comment to a photo and sends a response containing the unique ID of the created comment.
// It performs the following actions:
// - Verifies if the user is authenticated and not banned by the photo owner.
// - Validates the length of the comment.
// - Inserts the comment into the database.
// - Returns the unique ID of the comment in the response.
//
// Possible HTTP Responses:
// - 201 Created: The comment has been successfully added.
// - 400 Bad Request: There is an issue with the request (e.g., invalid data, comment too long).
// - 403 Forbidden: The user is not authorized to perform this action.
// - 500 Internal Server Error: An error occurred while processing the request.
func (rt *_router) postComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	photoOwnerId := ps.ByName("id")
	requestingUserId := extractBearer(r.Header.Get("Authorization"))

	if isNotLogged(requestingUserId) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Check if the requesting user is banned by the photo owner
	banned, err := rt.db.BannedUserCheck(
		User{IdUser: requestingUserId}.ToDatabase(),
		User{IdUser: photoOwnerId}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("post-comment/db.BannedUserCheck: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		w.WriteHeader(http.StatusForbidden) //user is banned by photo owner
		return
	}

	// decodeComment decodes the JSON request body into a Comment struct.
	var comment Comment
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("post-comment/Decode: failed to decode request body json")
		return
	}

	// validateCommentLength checks if the comment length is within the allowed limit.
	// If the length exceeds the limit, it returns an error indicating the comment is too long.
	if len(comment.Comment) > 30 {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("post-comment: comment length exceeds 30 characters")
		return
	}

	// parsePhotoID converts the photo ID from a string to an int64.
	// If the conversion fails, it returns an error indicating the photo ID is invalid.
	photoID, err := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("post-comment/ParseInt: failed to parse photo_id to int64")
		return
	}

	// Function call to db for comment creation
	commentId, err := rt.db.CommentPhoto(
		PhotoId{IdPhoto: photoID}.ToDatabase(),
		User{IdUser: requestingUserId}.ToDatabase(),
		comment.ToDatabase())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("post-comment/db.CommentPhoto: failed to insert comment")
		return
	}

	w.WriteHeader(http.StatusCreated)

	// The response body will contain the unique id of the comment
	err = json.NewEncoder(w).Encode(CommentId{IdComment: commentId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("post-comment/Encode: failed convert photo_id to int64")
		return
	}
}
