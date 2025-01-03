// Package jfif implements writing JPEG images with fixed JFIF header.
package jfif

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"io"
)

func GrayImage(m image.Image) *image.Gray {
	bounds := m.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	grayImg := image.NewGray(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Get the color of the pixel at (x, y)
			r, g, b, _ := m.At(x, y).RGBA()

			// Convert the color to grayscale using the luminance formula
			// Luminance = 0.299*R + 0.587*G + 0.114*B
			grayValue := uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256)
			// Set the grayscale pixel value
			grayImg.Set(x, y, color.Gray{Y: grayValue})
		}
	}
	return grayImg
}

var naiveJFIFHeader = []byte{
	0xFF, 0xD8, // SOI
	0xFF, 0xE0, // APP0 Marker
	0x00, 0x10, // Length
	0x4A, 0x46, 0x49, 0x46, 0x00, // JFIF\0
	0x01, 0x02, // 1.02
	0x00,       // Density type
	0x00, 0x01, // X Density
	0x00, 0x01, // Y Density
	0x00, 0x00, // No Thumbnail
}

// Encode writes the Image m to w in JFIF 1.02 compatible format with
// the given options. The JFIF header cannot be configured.
func Encode(w io.Writer, m image.Image, o *jpeg.Options) error {
	// Convert the image to grayscale with 256 levels (8 bits)
	grayImg := GrayImage(m)	
	buf := bytes.NewBuffer(nil)
	err := jpeg.Encode(buf, grayImg, o)
	if err != nil {
		return err
	}

	// Connect header and body
	body := buf.Bytes()[2:]
	_, err = w.Write(naiveJFIFHeader)
	if err != nil {
		return err
	}
	_, err = w.Write(body)
	if err != nil {
		return err
	}

	return nil
}
