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
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/img", imgHandler)

	log.Println("Starting server...")
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("cannot start, need a PORT")
		os.Exit(1)
	}
	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(indexHTML))
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

	log.Printf("url=\"%s\" size=%s get=%s resize=%s inbytes=%d outbytes=%d\n",
		imgURL, size, getDuration, resizeDuration, len(imgBytes), len(newImgBytes))

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

const indexHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
<title>gomagick</title>
</head>

<body>

<h1>gomagick</h1>
<p>
   Image resizing server written in Go based on ImageMagick, proof of concept.
   <br/>
   <a href="https://github.com/mat/gomagick">https://github.com/mat/gomagick</a>
</p>

<h2>Usage</h2>

<pre>http://gomagick.herokuapp.com/img?<strong>url</strong>=...&<strong>size</strong>=500x500</pre>

<h2>Examples</h2>

<table>
   <tr>
      <th>size param</th>
      <th></th>
   </tr>
   <tr>
      <td>100%</td>
      <td><img src="/img?size=100%25&amp;url=https://github.com/apple-touch-icon.png"></td>
   </tr>
   <tr>
      <td>50%</td>
      <td><img src="/img?size=50%25&amp;url=https://github.com/apple-touch-icon.png"></td>
   </tr>
   <tr>
      <td>150%</td>
      <td><img src="/img?size=150%25&amp;url=https://github.com/apple-touch-icon.png"></td>
   </tr>
   <tr>
      <td>100x100</td>
      <td><img src="/img?size=100x100&amp;url=https://github.com/apple-touch-icon.png"></td>
   </tr>
   <tr>
      <td>100x200!</td>
      <td><img src="/img?size=100x200!&amp;url=https://github.com/apple-touch-icon.png"></td>
   </tr>
</table>

</body>
</html>
`
