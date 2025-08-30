package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/image/webp"
)

func isValidType(ext string) bool {
	switch ext {
	case "png", "jpeg", "jpg", "webp":
		return true
	default:
		return false
	}
}

func decode(ext string, f io.Reader) (image.Image, error) {
	switch ext {
	case ".webp":
		img, err := webp.Decode(f)
		if err != nil {
			return nil, err
		}
		return img, err
	case ".jpg", ".jpeg":
		img, err := jpeg.Decode(f)
		if err != nil {
			return nil, err
		}
		return img, err
	case ".png":
		img, err := png.Decode(f)
		if err != nil {
			return nil, err
		}
		return img, err
	default:
		return nil, fmt.Errorf("invalid file format")
	}
}

func Encode(path, name, ext string, img image.Image) error {
	switch ext {
	case "png":
		pngFile, err := os.Create(path + "/" + name + "." + ext)
		if err != nil {
			return err
		}
		defer pngFile.Close()

		if err := png.Encode(pngFile, img); err != nil {
			return err
		}
		return nil

	case "jpg", "jpeg":
		jpegFile, err := os.Create(path + "." + ext)
		if err != nil {
			return err
		}
		defer jpegFile.Close()

		if err := jpeg.Encode(jpegFile, img, &jpeg.Options{Quality: 90}); err != nil {
			return err
		}
		return nil

	default:
		return fmt.Errorf("invalid file format")
	}
}

func ConvertFormat(url, newExt string) (string, error) {
	if !isValidType(newExt) {
		return "", fmt.Errorf("invalid file format")
	}

	var err error
	output := filepath.Base(url)
	oldExt := filepath.Ext(output)
	output = output[0 : len(output)-len(filepath.Ext(output))]

	// open the input file
	f, err := os.Open(url)
	if err != nil {
		return "", fmt.Errorf("unable to open file for reading")
	}

	img, err := decode(oldExt, f)
	if err != nil {
		return "", err
	}

	if err := Encode(url[0:len(url)-len(output+"."+oldExt)], output, newExt, img); err != nil {
		return "", err
	}

	return output, nil
}
