package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"wasaphoto-1764024/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// // deletePhoto handles the deletion of a photo, including removing the photo's file
// // and associated data (likes, comments) from the database.
// //
// // HTTP Responses:
// // - 204 No Content: Photo successfully deleted.
// // - 401 Unauthorized: The user is not authorized to delete the photo.
// // - 500 Internal Server Error: There was an issue with deleting the photo or file system operations.
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Extract the bearer token to authenticate the user

	bearerAuth := extractBearer(r.Header.Get("Authorization"))
	photoIdStr := ps.ByName("photo_id")

	// Check if the requesting user is authorized to delete the photo
	valid := validateRequestingUser(ps.ByName("id"), bearerAuth)
	if valid != 0 {
		ctx.Logger.Warnf("Unauthorized access attempt by user %s to delete photo", bearerAuth)
		w.WriteHeader(valid)
		return
	}

	// Convert the photo ID from string to int64
	photoInt, err := strconv.ParseInt(photoIdStr, 10, 64)
	if err != nil {
		ctx.Logger.WithError(err).Errorf("Invalid photo ID format: %s", photoIdStr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Call the database function to remove the photo and its associated data
	err = rt.db.RemovePhoto(
		User{IdUser: bearerAuth}.ToDatabase(),
		PhotoId{IdPhoto: photoInt}.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Errorf("Database error while removing photo with ID %d", photoInt)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the path to the user's photo folder
	pathPhoto, err := getUserPhotoFolder(bearerAuth)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error retrieving user's photo folder")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Full path to the photo file
	photoPath := filepath.Join(pathPhoto, photoIdStr)

	// Remove the photo file from the file system
	err = os.Remove(photoPath)
	if err != nil {
		// Log error if file doesn't exist but don't return an error to maintain idempotency
		ctx.Logger.WithError(err).Errorf("Failed to remove photo file: %s", photoPath)
	}

	// Return 204 No Content to indicate successful deletion
	w.WriteHeader(http.StatusNoContent)
}
