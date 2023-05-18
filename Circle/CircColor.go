package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		drawRandomColorFlash(w)
	})

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func drawRandomColorFlash(w http.ResponseWriter) {
	width := 1920
	height := 937

	rand.Seed(time.Now().UnixNano())

	// Create a new RGBA image with white background
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Rect.Canon(), &image.Uniform{color.Gray{}}, image.ZP, draw.Src)

	// Generate random circle properties
	radius := rand.Intn(100) + 50
	circleColor := randomColor()

	// Calculate the center of the image
	centerX := width / 2
	centerY := height / 2

	// Draw a circle in the image
	drawCircle(img, centerX, centerY, radius, circleColor)

	// Encode the image and write the response
	w.Header().Set("Content-Type", "image/png")
	err := png.Encode(w, img)
	if err != nil {
		log.Println("Failed to encode image:", err)
	}
}

func drawCircle(img draw.Image, centerX, centerY, radius int, color color.RGBA) {
	radiusSquared := radius * radius

	bounds := img.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			dx := x - centerX
			dy := y - centerY
			distanceSquared := dx*dx + dy*dy

			if distanceSquared <= radiusSquared {
				img.Set(x, y, color)
			}
		}
	}
}

func randomColor() color.RGBA {
	return color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}
}
