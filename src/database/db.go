package database

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"time"

	"github.com/disintegration/imaging"
	"github.com/patrickmn/go-cache"
)

var DBConnection *sql.DB
var Cache *cache.Cache

func StartDBAndCache() {
	// Spin up the database connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DBConnection, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = DBConnection.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Database Ready")

	Cache = cache.New(5*time.Minute, 10*time.Minute)

	fmt.Println("Cache Ready")
}

func CreateThumbnails() {
	fmt.Println("Creating Thumbnails")
	const existingImageIDsCacheQuery = `SELECT id FROM pictures
	UNION
	SELECT id FROM build_log_picture
	UNION
	SELECT ident AS id FROM plane
	UNION
	SELECT id FROM build`

	rows, _ := DBConnection.Query(existingImageIDsCacheQuery)
	for rows.Next() {
		var imageId string
		rows.Scan(&imageId)

		GetImage(imageId, true)
	}
}

func GetImage(imageId string, thumbnail bool) ([]byte, error) {
	// Check cache first if thumbnail
	if thumbnail {
		cachedThumbnail, found := Cache.Get(imageId)
		if found {
			return cachedThumbnail.([]byte), nil
		}
	}

	const imageQuery = `
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
	err := DBConnection.QueryRow(imageQuery, imageId).Scan(&imageData)
	if err != nil {
		var emptyBytes []byte
		return emptyBytes, err
	}

	if !thumbnail {
		// Resize to 800 height (to match modal on UI)
		// image.Image from bytes
		img, _, _ := image.Decode(bytes.NewReader(imageData))
		// Resize
		dstImage := imaging.Resize(img, 0, 800, imaging.Lanczos)
		// Back to []byte
		buf := new(bytes.Buffer)
		jpeg.Encode(buf, dstImage, nil)
		bytes := buf.Bytes()

		return bytes, nil
	}

	if len(imageData) == 0 {
		var emptyBytes []byte
		err := errors.New("image is null")
		return emptyBytes, err
	}

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
	Cache.Set(imageId, bytes, cache.NoExpiration)

	fmt.Printf("Set cache for %s\n", imageId)

	return bytes, nil
}
