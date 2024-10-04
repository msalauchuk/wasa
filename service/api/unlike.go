package api

import (
	"net/http"
	"strconv"
	"wasaphoto-1764024/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// unlike removes a like from a photo by a user.

// 1. Validate the requesting user's authentication using the bearer token.
// 2. Prevent self-unlike, as users cannot like or unlike their own photos (returns 400 Bad Request).
// 3. Check if the requesting user is banned by the photo author (returns 403 Forbidden if true).
// 4. Parse the photo ID from the request path and verify it is valid.
// 5. Remove the like from the photo in the database.
// 6. Respond with 204 No Content to indicate a successful operation.
//
// Response Codes:
// - 204 No Content: Successfully removed the like from the photo.
// - 400 Bad Request: The user attempted to unlike their own photo, or an invalid photo ID was provided.
// - 403 Forbidden: The requesting user has been banned by the photo author or is not logged in.
// - 500 Internal Server Error: An error occurred while processing the request in the database.

func (rt *_router) unlike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Extract Bearer Token
	requestingUserId := extractBearer(r.Header.Get("Authorization"))

	// Check if the user is logged in
	if isNotLogged(requestingUserId) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Extract path parameters
	photoAuthor := ps.ByName("id")
	photoIdStr := ps.ByName("photo_id")

	// Prevent self-unlike: users cannot unlike their own photo
	if photoAuthor == requestingUserId {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Check if the requesting user wasn't banned by the photo owner
	banned, err := rt.db.BannedUserCheck(
		User{IdUser: requestingUserId}.ToDatabase(),
		User{IdUser: photoAuthor}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("unlike/db.BannedUserCheck: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		// User was banned by owner, can't perform unlike action
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Parse the photo ID from string to int64
	photoId, err := strconv.ParseInt(photoIdStr, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("remove-like/ParseInt: error converting photo_id to int64")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Remove the like in the database via db function
	err = rt.db.UnlikePhoto(
		PhotoId{IdPhoto: photoId}.ToDatabase(),
		User{IdUser: requestingUserId}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("remove-like/db.UnlikePhoto: error executing delete query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 204 No Content to indicate successful deletion
	w.WriteHeader(http.StatusNoContent)
}
