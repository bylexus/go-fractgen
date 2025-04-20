package lib

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
)

type FractImage struct {
	*image.RGBA
	fractResults []FractFunctionResult
}

func NewFractImage(width, height int) *FractImage {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	return &FractImage{img, make([]FractFunctionResult, width*height)}
}

func (img *FractImage) EncodePng(w io.Writer) error {
	return png.Encode(w, img)
}

func (img *FractImage) EncodeJpeg(w io.Writer) error {
	return jpeg.Encode(w, img, &jpeg.Options{Quality: 90})
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
