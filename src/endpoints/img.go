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

func Img(w http.ResponseWriter, r *http.Request) {
	var thumbnail bool

	thumbnailQ, ok := r.URL.Query()["thumb"]
	if !ok || len(thumbnailQ) < 1 {
		thumbnail = false
	} else if thumbnailQ[0] == "true" {
		thumbnail = true
	}

	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids) < 1 {
		fmt.Fprintf(w, "id is required")
		return
	}
	id := ids[0]

	// If thumbnail is true, check cache and return if set, skip everything below
	if thumbnail {
		// Get thumbnail image from cache if it exists
		thumbnail, found := database.Cache.Get(id)
		if found {
			w.Header().Set("Content-Type", "image/jpeg")
			w.Write(thumbnail.([]byte))

			return
		}
	}

	imageQuery := `
		SELECT data FROM (
			SELECT id, data FROM pictures
			UNION
			SELECT id, data FROM build_log_picture
			UNION
			SELECT ident AS id, pic AS data FROM plane
			UNION
			SELECT id, cover AS data FROM build
		) as allimages
		WHERE id = $1
	`

	var imageData []byte
	err := database.DBConnection.QueryRow(imageQuery, id).Scan(&imageData)
	if err != nil {
		fmt.Fprintf(w, "Not Found")
		fmt.Println(err)
		return
	}

	if !thumbnail {
		// If thumbnail is false, just send the image as is
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(imageData)
	} else {
		// Create thumbnail image
		// image.Image from bytes
		img, _, _ := image.Decode(bytes.NewReader(imageData))

		// Resize
		dstImage := imaging.Resize(img, 0, 200, imaging.Lanczos)
		// Back to []byte
		buf := new(bytes.Buffer)
		jpeg.Encode(buf, dstImage, nil)
		bytes := buf.Bytes()

		// Save to cache for next time
		database.Cache.Set(id, bytes, cache.NoExpiration)

		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(bytes)
	}
}
