package main

import (
	"bytes"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mat/magick"

	// Supported image formats:
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func main() {
	http.HandleFunc("/", imgHandler)

	log.Println("Starting server...")
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("cannot start, need a PORT")
		os.Exit(1)
	}
	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func imgHandler(w http.ResponseWriter, r *http.Request) {
	imgURL := r.FormValue("url")
	if imgURL == "" {
		http.Error(w, "missing parameter: url (to image)", 401)
		return
	}
	size := r.FormValue("size")
	if size == "" {
		http.Error(w, "missing parameter: size (300x300, 50%)", 401)
		return
	}

	getStarted := time.Now()
	imgBytes, e := getImageData(imgURL)
	getDuration := time.Since(getStarted)
	if e != nil {
		httpError(w, e)
		return
	}

	imgFormat, e := detectFormat(imgBytes)
	if e != nil {
		httpError(w, e)
		return
	}

	resizeStarted := time.Now()
	newImgBytes, e := resizeImage(imgBytes, imgFormat, size)
	resizeDuration := time.Since(resizeStarted)
	if e != nil {
		httpError(w, e)
		return
	}

	log.Printf("get=%s resize=%s\n", getDuration, resizeDuration)

	w.Header().Set(contentType, fmt.Sprintf("image/%s", imgFormat))
	w.Header().Set(xTimings, fmt.Sprintf("get=%v, resize=%v", getDuration, resizeDuration))

	w.Write(newImgBytes)
}

func getImageData(imgURL string) ([]byte, *serverError) {
	imgReponse, err := http.Get(imgURL)
	if err != nil {
		return nil, &serverError{fmt.Sprintf("could not fetch image: %s", err), 400}
	}

	imgBytes, err := ioutil.ReadAll(imgReponse.Body)
	if err != nil {
		return nil, &serverError{fmt.Sprintf("could read image data: %s", err), 400}
	}
	defer imgReponse.Body.Close()

	return imgBytes, nil
}

func detectFormat(imgBytes []byte) (string, *serverError) {
	_, imgFormat, err := image.DecodeConfig(bytes.NewReader(imgBytes))
	if err != nil {
		return "", &serverError{fmt.Sprintf("could not detect image format: %s", err), 501}
	}
	return imgFormat, nil
}

func resizeImage(imgBytes []byte, imgFormat string, size string) ([]byte, *serverError) {
	img, err := magick.NewFromBlob(imgBytes, imgFormat)
	defer img.Destroy()
	if err != nil {
		return nil, &serverError{fmt.Sprintf("init failed: %s", err), 501}
	}

	err = img.Resize(size)
	if err != nil {
		return nil, &serverError{fmt.Sprintf("resize failed: %s", err), 501}
	}

	newImgBytes, err := img.ToBlob(imgFormat)
	if err != nil {
		return nil, &serverError{fmt.Sprintf("encoding failed: %s", err), 501}
	}

	return newImgBytes, nil
}

func httpError(w http.ResponseWriter, e *serverError) {
	http.Error(w, e.message, e.status)
}

type serverError struct {
	message string
	status  int
}

const contentType = "Content-Type"
const xTimings = "X-Image-Timings"
