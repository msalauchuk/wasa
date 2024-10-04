package api

import (
	"net/http"
	"wasaphoto-1764024/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// unban removes a user from the banned list of another user.
//

// This function first validates the requesting user's identity. If valid, it checks whether the user is trying to unban themselves,
// which is not allowed. If the operation is valid, the user is unbanned in the database. Otherwise, it responds with an appropriate
// error status.
//
// Response Codes:
// - 204 No Content: Successfully unbanned the user.
// - 400 Bad Request: Attempt to unban oneself or other invalid input.
// - 401 Unauthorized: The requesting user is not authorized to perform this action.
// - 500 Internal Server Error: An error occurred while processing the database request.
func (rt *_router) unban(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	bearerToken := extractBearer(r.Header.Get("Authorization"))

	// Extract IDs from path parameters
	pathId := ps.ByName("id")
	userToUnban := ps.ByName("banned_id")

	// Validate that the requesting user is the same as the path user
	valid := validateRequestingUser(pathId, bearerToken)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Users cannot unban themselves StatusBadRequest
	if userToUnban == bearerToken {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Remove the ban from the database
	err := rt.db.UnbanUser(
		User{IdUser: pathId}.ToDatabase(),
		User{IdUser: userToUnban}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("unban: error executing delete query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with 204 No Content to indicate success without response body
	w.WriteHeader(http.StatusNoContent)
}
