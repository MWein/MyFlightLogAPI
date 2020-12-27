package endpoints

import (
	"fmt"
	"net/http"

	"github.com/MWein/MyFlightLogAPI/src/database"
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

	image, err := database.GetImage(id, thumbnail)

	if err != nil {
		fmt.Fprintf(w, "Not Found")
		fmt.Println(err)
		return
	} else {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(image)
	}
}
