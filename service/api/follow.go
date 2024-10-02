package api

import (
	"net/http"
	"wasaphoto-1764024/service/api/reqcontext"
	"wasaphoto-1764024/service/database"

	"github.com/julienschmidt/httprouter"
)

// follow adds the requesting user to the followers list of another user.
// The user cannot follow themselves, and if the user was banned by the target user, the follow action is forbidden.
func (rt *_router) follow(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Extract IDs from the request path and authorization header
	targetUserId := ps.ByName("id")                                  // ID of the user to be followed
	requestingUserId := extractBearer(r.Header.Get("Authorization")) // ID of the authenticated user

	// Users cannot follow themselves
	if requestingUserId == targetUserId {
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request: Can't follow oneself
		return
	}

	// Check if the id of the follower in the request is the same of the bearer and the path parameter
	if ps.ByName("follower_id") != requestingUserId {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the requesting user was banned by the target user
	banned, err := rt.db.BannedUserCheck(
		database.User{IdUser: requestingUserId},
		database.User{IdUser: targetUserId})
	if err != nil {
		ctx.Logger.WithError(err).Error("follow/db.BannedUserCheck: error executing query")
		w.WriteHeader(http.StatusInternalServerError) // 500 Internal Server Error
		return
	}
	// If the requesting user is banned, the follow action is forbidden
	if banned {
		w.WriteHeader(http.StatusForbidden) // 403 Forbidden: User is banned
		return
	}

	// Add the new follower to the database
	err = rt.db.FollowUser(
		User{IdUser: requestingUserId}.ToDatabase(),
		User{IdUser: targetUserId}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("follow/db.FollowUser: error executing insert query")
		w.WriteHeader(http.StatusInternalServerError) // 500 Internal Server Error
		return
	}

	// Respond with 204 No Content: Follow operation successful with no content to return
	w.WriteHeader(http.StatusNoContent)
}
