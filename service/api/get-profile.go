package api

import (
	"encoding/json"
	"net/http"
	"wasaphoto-1764024/service/api/reqcontext"
	"wasaphoto-1764024/service/database"

	"github.com/julienschmidt/httprouter"
)

// getUserProfile handles the retrieval of complete profile information for a requested user.
// It performs several checks including bans between users, and returns profile data like followers,
// following users, and photos.
//
// Authorization: The requesting user must provide a valid Bearer token.
// The requested user's profile will only be fully visible if the requesting user is not banned by the profile owner
// and vice versa.
//
// HTTP Responses:
// - 200 OK: Complete profile retrieved successfully.
// - 204 No Content: Requested user does not exist.
// - 206 Partial Content: Requested profile is partially available due to bans.
// - 403 Forbidden: Requesting user is banned by the profile owner.
// - 500 Internal Server Error: Any database or internal error.
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	requestingUserId := extractBearer(r.Header.Get("Authorization"))
	requestedUser := ps.ByName("id")

	// Initialize empty slices for followers, following, and photos
	var followers, following []database.User
	var photos []database.Photo

	// Check if the requesting user is banned by the requested profile owner
	userBanned, err := rt.db.BannedUserCheck(User{IdUser: requestingUserId}.ToDatabase(),
		User{IdUser: requestedUser}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile/db.BannedUserCheck/userBanned: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userBanned {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Check if the requested profile was banned by the requesting user. If it's true respond with partial content
	requestedProfileBanned, err := rt.db.BannedUserCheck(User{IdUser: requestedUser}.ToDatabase(),
		User{IdUser: requestingUserId}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile/db.BannedUserCheck/requestedProfileBanned: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if requestedProfileBanned {
		w.WriteHeader(http.StatusPartialContent)
		return
	}

	// Verify if the requested user exists in the system
	userExists, err := rt.db.CheckUser(User{IdUser: requestedUser}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile: error checking user existence")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !userExists {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Fetch the user's followers
	followers, err = rt.db.GetFollowers(User{IdUser: requestedUser}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile: error retrieving followers")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Fetch the users the requested user is following
	following, err = rt.db.GetFollowing(User{IdUser: requestedUser}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile: error retrieving following")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Fetch the user's nickname

	nickname, err := rt.db.GetNickname(User{IdUser: requestedUser}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile/db.GetNickname: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Fetch the requested user's photos visible to the requesting user

	photos, err = rt.db.GetPhotosList(User{IdUser: requestingUserId}.ToDatabase(), User{IdUser: requestedUser}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("getUserProfile/db.GetPhotosList: error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return the complete profile as a JSON response
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(CompleteProfile{
		Name:      requestedUser,
		Nickname:  nickname,
		Followers: followers,
		Following: following,
		Posts:     photos,
	})

}
