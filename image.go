package framework

import (
	"errors"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

type Image struct {
	img image.Image
	ext string
}

func(img *Image) Set(path string) error {
	ext := strings.ToLower(filepath.Ext(path))
	img.ext = ext
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	switch ext {
	case ".png" :
		img.img, err = png.Decode(file)
	case ".gif" :
		img.img, err = gif.Decode(file)
	default:
		img.img, err = jpeg.Decode(file)
	}
	if err != nil {
		return err
	}
	file.Close()
	return  nil
}

func(img *Image) Resize(newSize uint,newPath string) error {
	if img.img == nil{
		return errors.New("Set Path Before")
	}
	image := resize.Resize(newSize, 0, img.img, resize.Lanczos3)
	out, err := os.Create(newPath)
	if err != nil {
		return err
	}
	defer out.Close()
	switch img.ext {
	case ".png" :
		png.Encode(out, image)
	case ".gif" :
		gif.Encode(out, image, nil)
	default:
		jpeg.Encode(out, image, nil)
	}
	return  nil
}