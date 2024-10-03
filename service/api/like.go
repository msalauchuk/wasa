package api

import (
	"net/http"
	"strconv"
	"wasaphoto-1764024/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// like allows an authenticated user to add a like to another user's photo.
// This operation is subject to several checks such as bans and consistency validation.
//
// HTTP Responses:
// - 204 No Content: Like successfully added.
// - 400 Bad Request: Request is invalid (e.g. trying to like own photo or inconsistent like ID).
// - 403 Forbidden: User is either not authenticated or is banned by the photo owner.
// - 500 Internal Server Error: A problem occurred when interacting with the database.
func (rt *_router) like(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Extract the photo author and requesting user's ID

	photoAuthor := ps.ByName("id")
	requestingUserId := extractBearer(r.Header.Get("Authorization"))
	pathLikeId := ps.ByName("like_id")

	// Check if the user is authenticated
	if isNotLogged(requestingUserId) {
		ctx.Logger.Warn("like: unauthenticated user tried to like a photo")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Prevent users from liking their own photos
	if photoAuthor == requestingUserId {
		ctx.Logger.Warn("like: user tried to like their own photo")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the requesting user has been banned by the photo's owner
	banned, err := rt.db.BannedUserCheck(
		User{IdUser: requestingUserId}.ToDatabase(),
		User{IdUser: photoAuthor}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("like: error executing BannedUserCheck query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		// If the user is banned by the photo owner, they cannot like the photo
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Follower id is not consistent with requesting user bearer token
	if pathLikeId != requestingUserId {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Convert the path parameters photo_id from string to int64
	photoId, err := strconv.ParseInt(ps.ByName("photo_id"), 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Error("like: error converting photo_id to int64")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Insert the like into the database
	err = rt.db.LikePhoto(
		PhotoId{IdPhoto: photoId}.ToDatabase(),
		User{IdUser: pathLikeId}.ToDatabase())
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Respond with 204 No Content to indicate success
	w.WriteHeader(http.StatusNoContent)
}
