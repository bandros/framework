package framework

import (
	"errors"
	"fmt"
	"image"
	"os"
)

func FileInfo(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func FileExist(path string) bool {
	_, err := FileInfo(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func IsDir(path string) bool {
	file, err := FileInfo(path)
	if err != nil {
		return false
	}
	return file.IsDir()
}

func GetSizeFile(path string) (int64, error) {
	file, err := FileInfo(path)
	if err != nil {
		return 0, err
	}
	return file.Size(), nil
}

func GetImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}
func RemoveFile(path string) error {
	if !FileExist(path) {
		return errors.New("File not exist")
	} else {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}
	return nil
}
