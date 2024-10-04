package api

import (
	"encoding/json"
	"net/http"
	"wasaphoto-1764024/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// putNickname handles the update of a user's nickname. It validates the user's identity,
// decodes the new nickname from the request body, and updates it in the database.

func (rt *_router) putNickname(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Extract user ID from the URL path
	pathId := ps.ByName("id")

	// Validate the user's identity for the operation
	valid := validateRequestingUser(pathId, extractBearer(r.Header.Get("Authorization")))
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Decode the new nickname from the request body
	var nickname Nickname
	err := json.NewDecoder(r.Body).Decode(&nickname)
	if err != nil {
		ctx.Logger.WithError(err).Error("update-nname: error decoding JSON request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Update the nickname in the database
	err = rt.db.ChangeNickname(
		User{IdUser: pathId}.ToDatabase(),
		nickname.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("update-nickname: error executing update query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with a 204 No Content status to indicate successful update
	w.WriteHeader(http.StatusNoContent)
}
