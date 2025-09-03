package transformations

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/disintegration/gift"
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

func clamp(val float64) float64 {
	if val > 255 {
		return 255
	}
	if val < 0 {
		return 0
	}
	return val
}

func openAndDecode(url, ext string) (image.Image, error) {
	f, err := os.Open(url)
	if err != nil {
		return nil, err
	}
	defer f.Close()

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

func ConvertFormat(url, newExt string) (string, error) {
	if !isValidType(newExt) {
		return "", fmt.Errorf("invalid file format")
	}

	var err error
	output := filepath.Base(url)
	oldExt := filepath.Ext(output)

	img, err := openAndDecode(url, oldExt)
	if err != nil {
		return "", err
	}

	if err := encode(SAVEDIR, getFileName(output, oldExt), newExt, img); err != nil {
		return "", err
	}

	return getFileName(output, newExt) + newExt, nil
}

func Resize(width, height int, url string) (string, error) {
	output, ext := conv(url)

	img, err := openAndDecode(url, ext)
	if err != nil {
		return "", err
	}

	resizedImg := imaging.Resize(img, width, height, imaging.Lanczos)
	if err := encode(SAVEDIR, getFileName(output, ext), ext, resizedImg); err != nil {
		return "", err
	}

	return output, nil
}

func Rotate(angle float64, url string) (string, error) {
	output := filepath.Base(url)
	ext := filepath.Ext(output)

	img, err := openAndDecode(url, ext)
	if err != nil {
		return "", err
	}

	rotatedImg := imaging.Rotate(img, angle, nil)
	if err := encode(SAVEDIR, getFileName(output, ext), ext, rotatedImg); err != nil {
		return "", err
	}

	return output, nil
}

func Crop(x1, y1, x2, y2 int, url string) (string, error) {
	output := filepath.Base(url)
	ext := filepath.Ext(output)

	img, err := openAndDecode(url, ext)
	if err != nil {
		return "", err
	}

	croppedImg := imaging.Crop(img, image.Rect(x1, y1, x2, y2))
	if err := encode(SAVEDIR, getFileName(output, ext), ext, croppedImg); err != nil {
		fmt.Println(err)
		fmt.Println("inside crop")
		return "", err
	}
	return output, nil
}

func Watermark(x, y int, url1, url2 string) (string, error) {
	watermarkImg := filepath.Base(url1)
	ext1 := filepath.Ext(watermarkImg)

	baseImg := filepath.Base(url2)
	ext2 := filepath.Ext(baseImg)

	watermark, err := openAndDecode(url1, ext1)
	if err != nil {
		return "", err
	}

	img, err := openAndDecode(url2, ext2)
	if err != nil {
		return "", err
	}

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

	if err := encode(SAVEDIR, getFileName(baseImg, ext2), ext2, m); err != nil {
		return "", err
	}
	return baseImg, nil
}

func Flip(direction string, url string) (string, error) {
	output := filepath.Base(url)
	ext := filepath.Ext(output)

	img, err := openAndDecode(url, ext)
	if err != nil {
		return "", err
	}

	// check  if it works otherwise try image.NewRGBA
	flippedImg := &image.NRGBA{}
	switch direction {
	case "horizontal":
		flippedImg = imaging.FlipH(img)
	case "vertical":
		flippedImg = imaging.FlipV(img)
	default:
		return "", fmt.Errorf("invalid flipping directon")
	}

	if err := encode(SAVEDIR, getFileName(output, ext), ext, flippedImg); err != nil {
		return "", err
	}

	return output, nil
}

func Mirror(url, direction string) (string, error) {
	output := filepath.Base(url)
	ext := filepath.Ext(url)

	img, err := openAndDecode(url, ext)
	if err != nil {
		return "", err
	}

	var dst *image.RGBA
	switch direction {
	case "horizontal":
		g := gift.New(gift.FlipHorizontal())
		dst = image.NewRGBA(g.Bounds(img.Bounds()))
		g.Draw(dst, img)
	case "vertical":
		g := gift.New(gift.FlipVertical())
		dst = image.NewRGBA(g.Bounds(img.Bounds()))
		g.Draw(dst, img)
	default:
		return "", fmt.Errorf("invalid flip direction")
	}

	if err := encode(SAVEDIR, getFileName(output, ext), ext, dst); err != nil {
		return "", err
	}

	return output, nil
}

func sepiaConv(img image.Image) *image.NRGBA {
	bounds := img.Bounds()
	sepiaImage := image.NewNRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			r, g, b, a := originalColor.RGBA()

			originalR := float64(r >> 8)
			originalG := float64(g >> 8)
			originalB := float64(b >> 8)

			newRed := 0.393*originalR + 0.769*originalG + 0.189*originalB
			newGreen := 0.349*originalR + 0.686*originalG + 0.168*originalB
			newBlue := 0.272*originalR + 0.534*originalG + 0.131*originalB

			newR := uint8(clamp(newRed))
			newG := uint8(clamp(newGreen))
			newB := uint8(clamp(newBlue))

			sepiaImage.Set(
				x,
				y,
				color.RGBA{R: newR, G: newG, B: newB, A: uint8(a >> 8)},
			)
		}
	}
	return sepiaImage
}

func Filters(filter, url string) (string, error) {
	output, ext := conv(url)

	img, err := openAndDecode(url, ext)
	if err != nil {
		return "", err
	}

	filterImg := &image.NRGBA{}
	switch filter {
	case "grayscale":
		filterImg = imaging.Grayscale(img)
	case "invert":
		filterImg = imaging.Invert(img)
	case "sepia":
		filterImg = sepiaConv(img)
	default:
		return "", fmt.Errorf("unsupported filter")
	}

	if err := encode(SAVEDIR, getFileName(output, ext), ext, filterImg); err != nil {
		return "", err
	}

	return output, nil
}
