package lib

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

type FractImage struct {
	*image.RGBA
}

func NewFractImage(width, height int) *FractImage {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	return &FractImage{img}
}

func (img *FractImage) SavePng(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func (img *FractImage) SaveJpeg(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := jpeg.Encode(f, img, &jpeg.Options{Quality: 90}); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}
	return nil
}
