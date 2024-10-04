package api

import (
	"net/http"
	"wasaphoto-1764024/service/api/reqcontext"
	"wasaphoto-1764024/service/database"

	"github.com/julienschmidt/httprouter"
)

// unfollow removes a user from the follower list of another user.
// 1. Validate the requesting user's identity by comparing the bearer token and the follower ID.
// 2. Prevent self-unfollow as it's an invalid action (returns 400 Bad Request).
// 3. Check if the requesting user is banned by the photo owner (returns 403 Forbidden if true).
// 4. If the user is allowed, remove the follower from the database.
// 5. Respond with 204 No Content to indicate a successful operation.
//
// Response Codes:
// - 204 No Content: Successfully unfollowed the user.
// - 400 Bad Request: Attempt to unfollow oneself.
// - 401 Unauthorized: The requesting user is not authorized to perform this action.
// - 403 Forbidden: The requesting user has been banned by the target user.
// - 500 Internal Server Error: An error occurred while processing the request in the database.

func (rt *_router) unfollow(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	requestingUserId := extractBearer(r.Header.Get("Authorization"))

	// Extract path parameters
	oldFollower := ps.ByName("follower_id")
	photoOwnerId := ps.ByName("id")

	// Validate that the follower ID matches the Bearer token to prevent impersonation
	valid := validateRequestingUser(oldFollower, requestingUserId)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Users cannot unfollow themselves; return a 400 Bad Request
	if photoOwnerId == requestingUserId {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Check if the requesting user has been banned by the photo owner
	banned, err := rt.db.BannedUserCheck(
		database.User{IdUser: requestingUserId},
		database.User{IdUser: photoOwnerId})
	if err != nil {
		ctx.Logger.WithError(err).Error("unfollow/rt.db.BannedUserCheck: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if banned {
		// User was banned, cannot perform the unfollow action
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Remove the follower from the database
	err = rt.db.UnfollowUser(
		User{IdUser: oldFollower}.ToDatabase(),
		User{IdUser: photoOwnerId}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("unfollow/rt.db.UnfollowUser: error executing delete query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 204 No Content to indicate success without response body
	w.WriteHeader(http.StatusNoContent)
}
