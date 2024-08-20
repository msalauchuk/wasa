package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"wasaphoto-1764024/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// sessionHandler manages the session creation for a user. It processes incoming user data from an HTTP request,
// performs validation, and interacts with the database to create a new user record. If the user creation is successful,
// it also initiates the setup of necessary file system directories for storing user-specific data (photos).
// This function ensures proper response codes and messages are sent back to the client depending on the outcome of
// these operations. Errors are logged and handled gracefully to provide clear feedback for troubleshooting.

func (rt *_router) sessionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	var user User
	// Decode the JSON body into the user struct; handle JSON parsing errors.
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		// Send a 400 Bad Request if the JSON is invalid.
		http.Error(w, `{"message": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	// Validate the user identifier; respond with 400 Bad Request if invalid.
	if !validIdentifier(user.IdUser) {
		http.Error(w, `{"message": "Invalid Identifier"}`, http.StatusBadRequest)
		return
	}

	// Check if the user already exists in the database.
	exists, err := rt.db.CheckUser(user.ToDatabase())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("session: failed to check if user exists")
		return
	}

	if exists {
		// If the user exists, return the existing user information
		existingUser, err := rt.db.GetUserByID(user.IdUser)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("session: failed to fetch existing user")
			return
		}

		// Return the existing user
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(existingUser); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("session: can't create response json")
		}
		return
	}

	// Attempt to create a new user record in the database.
	if err := rt.db.CreateUser(user.ToDatabase()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("session: failed to create user in database")
		return
	}

	if err := createUserFolder(user.IdUser, ctx); err != nil {
		ctx.Logger.WithError(err).Error("session: can't create user's photo folder")
		http.Error(w, `{"message": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	// Respond with 201 Created and encode the user object into the response body.
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		// Log and handle JSON encoding errors by responding with 500 Internal Server Error.
		ctx.Logger.WithError(err).Error("session: can't create response json")
		http.Error(w, `{"message": "Internal Server Error"}`, http.StatusInternalServerError)
	}
}

// Function that creates a new subdir for the specified user
func createUserFolder(identifier string, ctx reqcontext.RequestContext) error {

	// Create the path media/useridentifier/ inside the project dir
	path := filepath.Join(photoFolder, identifier)

	// To the previously created path add the "photos" subdir
	err := os.MkdirAll(filepath.Join(path, "photos"), os.ModePerm)
	if err != nil {
		ctx.Logger.WithError(err).Error("session/createUserFolder:: error creating directories for user")
		return err
	}
	return nil
}
