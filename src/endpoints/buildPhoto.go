package endpoints

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"

	"github.com/MWein/MyFlightLogAPI/src/database"
	"github.com/disintegration/imaging"
	"github.com/patrickmn/go-cache"
)

func BuildPhoto(w http.ResponseWriter, r *http.Request) {
	id, ok := r.URL.Query()["imgid"]

	if !ok || len(id[0]) < 1 {
		fmt.Fprintf(w, "Image ID is required")
		return
	}

	// Get thumbnail image from cache if it exists
	thumbnail, found := database.Cache.Get(id[0])
	if found {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(thumbnail.([]byte))

		return
	}

	// If no thumbnail is found in the cache, gotta make one

	var imageData []byte
	err := database.DBConnection.QueryRow("SELECT data FROM build_log_picture WHERE id = $1", id[0]).Scan(&imageData)
	if err != nil {
		fmt.Fprintf(w, "Not Found")
		return
	}

	// Create thumbnail image
	// image.Image from bytes
	img, _, _ := image.Decode(bytes.NewReader(imageData))

	// Resize
	dstImage := imaging.Resize(img, 0, 200, imaging.Linear)
	// Back to []byte
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, dstImage, nil)
	bytes := buf.Bytes()

	database.Cache.Set(id[0], bytes, cache.NoExpiration)

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(bytes)
}
