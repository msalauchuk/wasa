package api

import (
	"net/http"
	"os"
	"path/filepath"
	"wasaphoto-1764024/service/api/reqcontext"

	"github.com/julienschmidt/httprouter"
)

// getPhoto serves the requested photo to the client.
// It constructs the file path based on the user's ID and photo ID from the URL parameters.
func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Costruisci il percorso del file
	photoPath := filepath.Join(photoFolder, ps.ByName("id"), "photos", ps.ByName("photo_id"))

	// Controlla se il file esiste
	if _, err := os.Stat(photoPath); os.IsNotExist(err) {
		ctx.Logger.WithError(err).Error("getPhoto: file not found")
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("getPhoto: error checking file existence")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, photoPath)

}
