package api

import (
	"net/http"
	"strconv"
	"wasaphoto-1764024/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// deleteComment removes a comment from a photo.
//   - The function checks if the user is logged in.
//   - It verifies if the user has been banned by the photo owner.
//   - It validates and converts the photo and comment IDs from strings to integers.
//   - Depending on the user's role (photo owner or comment author), it deletes the comment.
func (rt *_router) deleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")
	requestingUserId := extractBearer(r.Header.Get("Authorization"))

	// Check if the user isn't logged
	if isNotLogged(requestingUserId) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Check if the requesting user wasn't banned by the photo owner
	photoOwnerId := ps.ByName("id")
	banned, err := rt.db.BannedUserCheck(
		User{IdUser: requestingUserId}.ToDatabase(),
		User{IdUser: photoOwnerId}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("post-comment/db.BannedUserCheck: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		w.WriteHeader(http.StatusForbidden) //User is banned by the photo owner
		return
	}

	// Convert the photo identifier from string to int64
	photoId, err := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("post-comment: failed convert photo_id to int64")
		return
	}

	// Convert the comment identifier from string to int64
	commentId, err := strconv.ParseInt(ps.ByName("comment_id"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("post-comment: failed convert photo_id to int64")
		return
	}

	// Remove the comment if the user is the author of the post
	if photoOwnerId == requestingUserId {

		err = rt.db.UncommentPhotoAuthor(
			PhotoId{IdPhoto: photoId}.ToDatabase(),
			CommentId{IdComment: commentId}.ToDatabase())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("post-comment: failed to execute query for comment removal by author")
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Remove the comment if the user is the author of the comment
	err = rt.db.UncommentPhoto(
		PhotoId{IdPhoto: photoId}.ToDatabase(),
		User{IdUser: requestingUserId}.ToDatabase(),
		CommentId{IdComment: commentId}.ToDatabase())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("post-comment: failed to execute query for comment removal")
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204
}
