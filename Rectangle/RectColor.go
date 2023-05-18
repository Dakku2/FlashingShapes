package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	imageWidth  = 1920
	imageHeight = 937
)

func main() {
	http.HandleFunc("/", colorRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func colorRequest(w http.ResponseWriter, r *http.Request) {
	// Generate a random color
	rand.Seed(time.Now().UnixNano())
	red := uint8(rand.Intn(256))
	green := uint8(rand.Intn(256))
	blue := uint8(rand.Intn(256))
	color := color.RGBA{red, green, blue, 255}

	// Create a new image and fill it with the random color
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{color}, image.ZP, draw.Src)

	// Encode the image as JPEG and write it as the response
	w.Header().Set("Content-Type", "image/jpeg")
	err := jpeg.Encode(w, img, nil)
	if err != nil {
		log.Println("Error encoding image:", err)
	}
}
