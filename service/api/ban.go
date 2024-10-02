package api

import (
	"net/http"
	"wasaphoto-1764024/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// ban adds a user to the banned list of another user.
// Only the owner of the account can perform this operation, and banning oneself is not allowed.
// The ban action also implies that the banned user and the banner no longer follow each other.
func (rt *_router) ban(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Extract IDs from the request path and authorization header
	ownerId := ps.ByName("id")
	toBandId := ps.ByName("banned_id")
	requestingUserId := extractBearer(r.Header.Get("Authorization"))

	// Check the user's identity for the operation (only owner of the account can add a banned user to that account list)
	if valid := validateRequestingUser(ownerId, requestingUserId); valid != 0 {
		w.WriteHeader(valid) // Return appropriate error code (e.g., 403 Forbidden)
		return
	}

	// Ensure that the user cannot ban themselves
	if requestingUserId == toBandId {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Add the banned user to the database
	if err := rt.db.BanUser(
		User{IdUser: ownerId}.ToDatabase(),
		User{IdUser: toBandId}.ToDatabase()); err != nil {
		ctx.Logger.WithError(err).Error("put-ban/db.BanUser: error executing insert query")
		w.WriteHeader(http.StatusInternalServerError) // Return 500 Internal Server Error on database failure
		return
	}

	// Ban implies removing the follow relationship between the requesting user and the banned user
	// First, unfollow the banned user if they are followed by the requesting user
	if err := rt.db.UnfollowUser(
		User{IdUser: requestingUserId}.ToDatabase(),
		User{IdUser: toBandId}.ToDatabase()); err != nil {
		ctx.Logger.WithError(err).Error("put-ban/db.UnfollowUser1: error executing insert query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Ensure the banned user no longer follows the requesting user
	if err := rt.db.UnfollowUser(
		User{IdUser: toBandId}.ToDatabase(),
		User{IdUser: requestingUserId}.ToDatabase()); err != nil {
		ctx.Logger.WithError(err).Error("put-ban/db.UnfollowUser2: error executing insert query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 204 http status
	w.WriteHeader(http.StatusNoContent)
}
