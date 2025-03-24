package main

import (
	"fmt"

	"github.com/bylexus/go-fract/lib"
)

func main() {

	// width := 1920
	// height := 1200
	presetsFile := "presets.json"

	colorPresets, fractalPresets := lib.ReadPresetJson(presetsFile)
	// fmt.Printf("Loaded %d color presets and %d fractal presets\n", len(colorPresets), len(fractalPresets))
	// fmt.Printf("Using fractal presets: %#v\n", fractalPresets)
	// fmt.Printf("Using color presets: %#v\n", colorPresets)

	// Standard mandelbrot view:
	// var fractal lib.Fractal = lib.NewMandelbrotFractal(width, height, -0.7, 0.0, 4.0, 50)
	// Seahorse Valley:
	// var fractal = lib.NewMandelbrotFractal(width, height, -0.87591, 0.20464, 0.53184, 100)
	// Fractal from preset:
	for i, p := range colorPresets {
		fractal, err := lib.NewFractalFromPresets(p, fractalPresets[6])
		if err != nil {
			panic(err)
		}
		img := fractal.CalcFractalImage()
		img.SavePng(fmt.Sprintf("image-color-preset-%02d.png", i))
	}
}
