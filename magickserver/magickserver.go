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

	// Needed for format detection:
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

func imgHandler(w http.ResponseWriter, r *http.Request) {
	imgURL := r.FormValue("url")
	if imgURL == "" {
		http.Error(w, "missing parameter: url", 401)
		return
	}

	size := r.FormValue("size")
	if size == "" {
		http.Error(w, "missing parameter: size", 401)
		return
	}

	getStarted := time.Now()
	imgReponse, err := http.Get(imgURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not fetch image: %s", err), 400)
		return
	}
	getDuration := time.Since(getStarted)

	imgBytes, err := ioutil.ReadAll(imgReponse.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not read image data: %s", err), 400)
		return
	}
	defer imgReponse.Body.Close()

	detectionStarted := time.Now()
	_, imgFormat, err := image.DecodeConfig(bytes.NewReader(imgBytes))
	if err != nil {
		http.Error(w, fmt.Sprintf("could not detect image format: %s", err), 501)
		return
	}
	detectionDuration := time.Since(detectionStarted)

	resizeStarted := time.Now()
	img, err := magick.NewFromBlob(imgBytes, imgFormat)
	defer img.Destroy()
	err = img.Resize(size)
	if err != nil {
		http.Error(w, fmt.Sprintf("resize failed: %s", err), 501)
		return
	}

	newImgBytes, err := img.ToBlob(imgFormat)
	if err != nil {
		http.Error(w, fmt.Sprintf("encoding failed: %s", err), 501)
		return
	}
	resizeDuration := time.Since(resizeStarted)
	log.Println("get=", getDuration, ", detection=", detectionDuration, ", resize=", resizeDuration)

	w.Header().Set(contentType, fmt.Sprintf("image/%s", imgFormat))
	w.Header().Set(xTimings, fmt.Sprintf("get=%v, resize=%v", getDuration, resizeDuration))

	w.Write(newImgBytes)
}

const contentType = "Content-Type"
const xTimings = "X-Image-Timings"

func main() {
	http.HandleFunc("/img", imgHandler)
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("cannot start, need a PORT")
		os.Exit(1)
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
