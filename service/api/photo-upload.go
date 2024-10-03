package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"wasaphoto-1764024/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// postPhoto manages the upload of a photo. It validates the format of the photo,
// saves it to the user's photo folder, and stores metadata in the database.
// - 201 Created: The photo has been successfully uploaded.
// - 400 Bad Request: The uploaded file is not in a valid format (jpeg/png).
// - 401 Unauthorized: The user is not authorized to upload the photo.
// - 500 Internal Server Error: An error occurred while processing the request.
func (rt *_router) postPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Set the content type for the response
	w.Header().Set("Content-Type", "application/json")

	auth := extractBearer(r.Header.Get("Authorization"))

	// Validate the requesting user's identity
	valid := validateRequestingUser(ps.ByName("id"), auth)
	if valid != 0 {
		w.WriteHeader(valid)
		return
	}

	// Initialize the photo metadata (with owner and upload date)
	photo := Photo{
		Owner: auth,
		Date:  time.Now().UTC(),
	}

	// Read the body content for validation and storage
	data, err := io.ReadAll(r.Body)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-upload: error reading body content")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Assign a new io.Reader to r.Body for further reading after validation
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	// Validate the photo format (jpeg or png)
	err = checkFormatPhoto(r.Body, io.NopCloser(bytes.NewBuffer(data)), ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("photo-upload: unsupported image format")
		// controllaerrore
		_ = json.NewEncoder(w).Encode(JSONErrorMsg{Message: IMG_FORMAT_ERROR_MSG})
		return
	}

	// Body has been read in the previous function,so it's necessary to reassign a io.ReadCloser to it
	r.Body = io.NopCloser(bytes.NewBuffer(data))

	// Store the photo metadata in the database and generate a unique ID for the photo
	photoIdInt, err := rt.db.CreatePhoto(photo.ToDatabase())
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-upload: database error while creating photo")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	photoId := strconv.FormatInt(photoIdInt, 10)

	// Retrieve the user's photo folder path or create it if it doesn't exist
	PhotoPath, err := getUserPhotoFolder(auth)
	if err != nil {
		ctx.Logger.WithError(err).Error("photo-upload: error accessing user's photo folder")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create an empty file for storing the body content (image)
	out, err := os.Create(filepath.Join(PhotoPath, photoId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("photo-upload: error creating local photo file")
		return
	}

	// Copy body content to the previously created file
	_, err = io.Copy(out, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("photo-upload: error copying body content into file photo")
		return
	}

	// Close the created file
	out.Close()

	// Respond with 201 Created and return the photo metadata
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(Photo{
		Comments: nil,
		Likes:    nil,
		Owner:    photo.Owner,
		Date:     photo.Date,
		PhotoId:  int(photoIdInt),
	})

}

// checkFormatPhoto validates whether the uploaded photo is in either jpeg or png format.
// It returns an error if the format is unsupported.
func checkFormatPhoto(body io.ReadCloser, newReader io.ReadCloser, ctx reqcontext.RequestContext) error {

	// Attempt to decode the image as a jpeg

	_, errJpg := jpeg.Decode(body)
	if errJpg != nil {
		body = newReader
		// Attempt to decode as png if jpeg decoding fails
		_, errPng := png.Decode(body)
		if errPng != nil {
			return errors.New(IMG_FORMAT_ERROR_MSG)
		}
		return nil
	}
	return nil
}

// getUserPhotoFolder returns the path to the user's photo folder.
// If the folder does not exist, it will be created.
// Example path: "./media/{user_id}/photos/"
func getUserPhotoFolder(userId string) (UserPhotoFoldrPath string, err error) {

	photoPath := filepath.Join(photoFolder, userId, "photos")

	return photoPath, nil
}
