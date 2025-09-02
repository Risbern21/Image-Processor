package transformations

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"golang.org/x/image/webp"
)

const SAVEDIR = "./dest/"

func getFileName(name, ext string) string {
	return name[:len(name)-len(ext)]
}

func conv(url string) (string, string) {
	return filepath.Base(url), filepath.Ext(url)
}

func isValidType(ext string) bool {
	switch ext {
	case ".png", ".jpeg", ".jpg", ".webp":
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

func encode(path, name, ext string, img image.Image) error {
	switch ext {
	case ".png":
		pngFile, err := os.Create(path + name + ext)
		if err != nil {
			return err
		}
		defer pngFile.Close()

		if err := png.Encode(pngFile, img); err != nil {
			return err
		}
		return nil

	case ".jpg", ".jpeg":
		jpegFile, err := os.Create(path + name + ext)
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

func ConvertFormat(url, newExt string) error {
	if !isValidType(newExt) {
		return fmt.Errorf("invalid file format")
	}

	var err error
	output := filepath.Base(url)
	oldExt := filepath.Ext(output)

	// open the input file
	f, err := os.Open(url)
	if err != nil {
		return fmt.Errorf("unable to open file for reading")
	}

	img, err := decode(oldExt, f)
	if err != nil {
		return err
	}

	if err := encode(SAVEDIR, getFileName(output, oldExt), newExt, img); err != nil {
		return err
	}

	return nil
}

func Resize(width, height int, url string) error {
	output, ext := conv(url)

	f, err := os.Open(url)
	if err != nil {
		return err
	}

	img, err := decode(ext, f)
	if err != nil {
		return err
	}
	defer f.Close()

	resizedImg := imaging.Resize(img, width, height, imaging.Lanczos)
	if err := encode(SAVEDIR, getFileName(output, ext)+"_resized", ext, resizedImg); err != nil {
		return err
	}

	return nil
}

func Rotate(angle float64, url string) error {
	output := filepath.Base(url)
	ext := filepath.Ext(output)

	f, err := os.Open(url)
	if err != nil {
		return err
	}

	img, err := decode(ext, f)
	if err != nil {
		return err
	}
	defer f.Close()

	rotatedImg := imaging.Rotate(img, angle, nil)
	if err := encode(SAVEDIR, getFileName(output, ext)+"_rotated", ext, rotatedImg); err != nil {
		return err
	}

	return nil
}

func Crop(x1, y1, x2, y2 int, url string) error {
	output := filepath.Base(url)
	ext := filepath.Ext(output)

	f, err := os.Open(url)
	if err != nil {
		return err
	}
	defer f.Close()

	img, err := decode(ext, f)
	if err != nil {
		return err
	}

	croppedImg := imaging.Crop(img, image.Rect(x1, y1, x2, y2))
	if err := encode(SAVEDIR, getFileName(output, ext), ext, croppedImg); err != nil {
		return err
	}
	return nil
}

func Watermark(x, y int, url1, url2 string) error {
	watermarkImg := filepath.Base(url1)
	ext1 := filepath.Ext(watermarkImg)

	baseImg := filepath.Base(url2)
	ext2 := filepath.Ext(baseImg)

	// open the watermark image and decode it
	wmb, err := os.Open(url1)
	if err != nil {
		return err
	}
	watermark, err := decode(ext1, wmb)
	if err != nil {
		return err
	}
	defer wmb.Close()

	// open the base image and decode it
	imgb, err := os.Open(url2)
	if err != nil {
		return err
	}
	img, err := decode(ext2, imgb)
	if err != nil {
		return err
	}
	defer imgb.Close()

	offset := image.Pt(x, y)
	b := img.Bounds()
	m := image.NewRGBA(b)
	draw.Draw(m, b, img, image.Point{0, 0}, draw.Src)
	draw.Draw(
		m,
		watermark.Bounds().Add(offset),
		watermark,
		image.Point{0, 0},
		draw.Over,
	)

	if err := encode(SAVEDIR, getFileName(baseImg, ext2)+"_watermarked", ext2, m); err != nil {
		return err
	}
	return nil
}

func Flip(direction string, url string) error {
	output := filepath.Base(url)
	ext := filepath.Ext(output)

	f, err := os.Open(url)
	if err != nil {
		return err
	}
	defer f.Close()

	img, err := decode(ext, f)
	if err != nil {
		return err
	}

	// check  if it works otherwise try image.NewRGBA
	flippedImg := &image.NRGBA{}
	switch direction {
	case "horizontal":
		flippedImg = imaging.FlipH(img)
	case "vertical":
		flippedImg = imaging.FlipV(img)
	default:
		return fmt.Errorf("invalid flipping directon")
	}

	if err := encode(SAVEDIR, getFileName(output, ext)+"_flip", ext, flippedImg); err != nil {
		return err
	}

	return nil
}

func MirrorHorizontal(url string) error {
	output := filepath.Base(url)
	ext := filepath.Ext(url)

	f, err := os.Open(url)
	if err != nil {
		return err
	}
	defer f.Close()

	img, err := decode(ext, f)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	newImg := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newImg.Set(x, y, newImg.At(width-1-x, y))
		}
	}

	if err := encode(SAVEDIR, getFileName(output, ext), ext, newImg); err != nil {
		return err
	}

	return nil
}

func MirrorVertical(url string) error {
	output := filepath.Base(url)
	ext := filepath.Ext(url)

	f, err := os.Open(url)
	if err != nil {
		return err
	}
	defer f.Close()

	img, err := decode(ext, f)
	if err != nil {
		return err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	newImg := image.NewRGBA(bounds)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			newImg.Set(x, y, newImg.At(x, height-1-y))
		}
	}

	if err := encode(SAVEDIR, getFileName(output, ext), ext, newImg); err != nil {
		return err
	}

	return nil
}

func Filters(filter, url string) error {
	output, ext := conv(url)

	f, err := os.Open(url)
	if err != nil {
		return err
	}
	defer f.Close()

	img, err := decode(ext, f)
	if err != nil {
		return err
	}

	filteImg := &image.NRGBA{}
	switch filter {
	case "grayscale":
		filteImg = imaging.Grayscale(img)
	case "invert":
		filteImg = imaging.Invert(img)
	default:
		return fmt.Errorf("unsupported filter")
	}

	if err := encode(SAVEDIR, getFileName(output, ext)+"filtered", ext, filteImg); err != nil {
		return err
	}

	return nil
}
