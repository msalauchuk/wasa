package api

import (
	"encoding/json"
	"net/http"
	"wasaphoto-1764024/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// searchUser retrieves all the users matching the query parameter and sends the response containing all the matches.
// The user must be authenticated (via Bearer token) to perform this operation.

// - 403 Forbidden: If the user is not logged in.
// - 500 Internal Server Error: If there is an issue querying the database.
// - 200 OK: Successfully retrieved users that match the query parameter.
func (rt *_router) searchUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Extract Bearer token from Authorization header
	identifier := extractBearer(r.Header.Get("Authorization"))

	// If the user is not logged in, respond with a 403 http status
	if identifier == "" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Extract the search query parameter from the URL
	searchID := r.URL.Query().Get("id")

	// Search for users in the database based on the provided query parameter
	res, err := rt.db.SearchUser(User{IdUser: identifier}.ToDatabase(), User{IdUser: searchID}.ToDatabase())
	if err != nil {
		// Log the error and respond with 500 Internal Server Error
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Database encountered an error during user search")
		_ = json.NewEncoder(w).Encode([]User{})
		return
	}

	// Set response status to 200 OK
	w.WriteHeader(http.StatusOK)

	// Encode the result as JSON and send the response
	if len(res) == 0 {
		_ = json.NewEncoder(w).Encode([]User{}) // send empty slice of Users
		return
	}
	_ = json.NewEncoder(w).Encode(res)
}
