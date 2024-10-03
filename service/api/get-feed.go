package api

import (
	"encoding/json"
	"net/http"
	"wasaphoto-1764024/service/api/reqcontext"
	"wasaphoto-1764024/service/database"

	"github.com/julienschmidt/httprouter"
)

// getFeed retrieves all the photos of the people that the user is following.
// It returns a JSON-encoded array of photos or an empty array if there are no photos to return.
func (rt *_router) getFeed(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Set the content type for the response

	w.Header().Set("Content-Type", "application/json")

	// Extract the user ID from the Authorization header (Bearer token)

	identifier := extractBearer(r.Header.Get("Authorization"))

	if identifier == "" {
		ctx.Logger.Error("getFeed: missing or invalid Authorization token")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Ensure the user is requesting their own feed
	if validateRequestingUser(ps.ByName("id"), identifier) != 0 {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Retrieve the list of users the requesting user is following
	followingUsers, err := rt.db.GetFollowing(User{IdUser: identifier}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getFeed: error retrieving following users")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var photos []database.Photo
	for _, followedUser := range followingUsers {

		userPhotos, err := rt.db.GetPhotosList(
			User{IdUser: identifier}.ToDatabase(),
			User{IdUser: followedUser.IdUser}.ToDatabase())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Add a limited number of photos from the followed user
		for i, photo := range userPhotos {
			if i >= database.PhotosPerUserHome {
				break
			}
			photos = append(photos, photo)
		}

	}
	// Respond with the collected photos json or an empty array if no photos are found
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(photos); err != nil {
		ctx.Logger.WithError(err).Error("getFeed: error encoding photos to JSON")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
