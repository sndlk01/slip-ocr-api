package ocr

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

func PreprocessImage(inputPath string) (string, error) {
	img, err := imaging.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}

	log.Printf("Original image size: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())

	img = imaging.Grayscale(img)
	log.Println("Converted to grayscale")

	img = imaging.AdjustContrast(img, 30)
	log.Println("Adjusted contrast")

	img = imaging.Sharpen(img, 2.0)
	log.Println("Sharpened image")

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	if width < 1000 || height < 1000 {
		if width > height {
			img = imaging.Resize(img, 1500, 0, imaging.Lanczos)
		} else {
			img = imaging.Resize(img, 0, 1500, imaging.Lanczos)
		}
		log.Printf("Resized image to: %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
	}

	outputPath := getPreprocessedPath(inputPath)
	err = imaging.Save(img, outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to save preprocessed image: %w", err)
	}

	log.Printf("Preprocessed image saved to: %s", outputPath)
	return outputPath, nil
}

func getPreprocessedPath(originalPath string) string {
	dir := filepath.Dir(originalPath)
	ext := filepath.Ext(originalPath)
	base := filepath.Base(originalPath)
	name := base[:len(base)-len(ext)]

	return filepath.Join(dir, name+"_processed"+ext)
}

// ConvertToJPEG converts an image to JPEG format if it's not already
func ConvertToJPEG(inputPath string) (string, error) {
	ext := filepath.Ext(inputPath)

	// If already JPEG, return as is
	if ext == ".jpg" || ext == ".jpeg" {
		return inputPath, nil
	}

	// Open the image
	file, err := os.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	var img image.Image

	// Decode based on extension
	switch ext {
	case ".png":
		img, err = png.Decode(file)
	default:
		img, _, err = image.Decode(file)
	}

	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Create output path
	outputPath := inputPath[:len(inputPath)-len(ext)] + ".jpg"

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Encode as JPEG
	err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 95})
	if err != nil {
		return "", fmt.Errorf("failed to encode JPEG: %w", err)
	}

	log.Printf("Converted image to JPEG: %s", outputPath)

	// Remove original file
	os.Remove(inputPath)

	return outputPath, nil
}
