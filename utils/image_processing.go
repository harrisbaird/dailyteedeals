package utils

import (
	"fmt"
	"image"
	"image/color"

	"github.com/disintegration/imaging"
)

var (
	cropSize = 30
	anchors  = []imaging.Anchor{
		imaging.TopLeft,
		imaging.TopRight,
		imaging.BottomLeft,
		imaging.BottomRight,
	}
)

type Transform struct {
	Size     int
	Path     string
	Format   imaging.Format
	MimeType string
}

func TransformImage(img image.Image, transform Transform, bgColor color.Color) image.Image {
	if transform.Size == 0 {
		return img
	}

	resizeWidth, resizeHeight := transform.Size, transform.Size
	bounds := img.Bounds()

	// Setting a dimension to zero resizes while maintaining aspect ratio.
	if bounds.Max.X > bounds.Max.Y {
		resizeHeight = 0
	} else {
		resizeWidth = 0
	}

	dstImage := imaging.Resize(img, resizeWidth, resizeHeight, imaging.CatmullRom)
	dstImage = imaging.Sharpen(dstImage, .5)

	// Create a new solid color image
	bgImage := imaging.New(transform.Size, transform.Size, bgColor)

	return imaging.OverlayCenter(bgImage, dstImage, 1)
}

func FindBackgroundColor(img image.Image) color.Color {
	histogram := map[color.Color]int{}
	var bestColor color.Color
	bestCount := 0

	for _, anchor := range anchors {
		cropped := imaging.CropAnchor(img, cropSize, cropSize, anchor)
		for pos := 0; pos < cropSize; pos++ {
			currColor := cropped.At(pos, pos)

			histogram[currColor]++
			currCount := histogram[currColor]
			if currCount > bestCount {
				bestCount = currCount
				bestColor = currColor
			}
		}
	}

	return bestColor
}

func ColorToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r/256, g/256, b/256)
}
