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

const (
	width  = 1920
	height = 937
)

// GenerateRandomColor generates a random RGB color.
func GenerateRandomColor() color.RGBA {
	return color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255, // Opaque
	}
}

// GenerateRandomTriangle generates a random colored triangle image.
func GenerateRandomTriangle() *image.RGBA {
	rand.Seed(time.Now().UnixNano())

	// Create a new RGBA image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	// Generate a random color
	color := GenerateRandomColor()

	// Calculate the triangle points
	x1, y1 := rand.Intn(width), rand.Intn(height)
	x2, y2 := rand.Intn(width), rand.Intn(height)
	x3, y3 := rand.Intn(width), rand.Intn(height)

	// Draw the triangle with the random color
	triangle := []image.Point{{x1, y1}, {x2, y2}, {x3, y3}}
	draw.DrawMask(img, img.Bounds(), &image.Uniform{color}, image.Point{}, &triangleImage{triangle}, image.Point{}, draw.Over)

	// Return the image
	return img
}

// triangleImage represents an image that draws a triangle.
type triangleImage struct {
	triangle []image.Point
}

// ColorModel returns the color model used by the triangle image.
func (t *triangleImage) ColorModel() color.Model {
	return color.AlphaModel
}

// Bounds returns the image bounds of the triangle image.
func (t *triangleImage) Bounds() image.Rectangle {
	minX, minY, maxX, maxY := t.triangle[0].X, t.triangle[0].Y, t.triangle[0].X, t.triangle[0].Y
	for _, p := range t.triangle {
		if p.X < minX {
			minX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	return image.Rect(minX, minY, maxX, maxY)
}

// At returns the color of the specified pixel in the triangle image.
func (t *triangleImage) At(x, y int) color.Color {
	p := image.Point{x, y}
	if pointInTriangle(p, t.triangle) {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}

// pointInTriangle checks if a given point is inside the triangle.
func pointInTriangle(p image.Point, triangle []image.Point) bool {
	p1, p2, p3 := triangle[0], triangle[1], triangle[2]
	d1 := sign(p, p1, p2)
	d2 := sign(p, p2, p3)
	d3 := sign(p, p3, p1)
	hasNeg := (d1 < 0) || (d2 < 0) || (d3 < 0)
	hasPos := (d1 > 0) || (d2 > 0) || (d3 > 0)
	return !(hasNeg && hasPos)
}

// sign calculates the sign of the cross product between vectors p1p2 and p1p3.
func sign(p1, p2, p3 image.Point) int {
	return (p1.X-p3.X)*(p2.Y-p3.Y) - (p2.X-p3.X)*(p1.Y-p3.Y)
}

// FlashTriangleHandler handles the HTTP request and renders the random colored triangle.
func FlashTriangleHandler(w http.ResponseWriter, r *http.Request) {
	// Generate a random colored triangle image
	img := GenerateRandomTriangle()

	// Set the Content-Type header to indicate it's an image
	w.Header().Set("Content-Type", "image/png")

	// Write the image to the response
	err := png.Encode(w, img)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", FlashTriangleHandler)
	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
