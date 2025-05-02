package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/bylexus/go-fract/cli"
	"github.com/bylexus/go-fract/lib"
)

func main() {
	cli := cli.Cli{}
	ctx := kong.Parse(&cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}

func demo() {
	presetsFile := "presets.json"

	presets := lib.ReadPresetJson(presetsFile)

	// Standard mandelbrot view:
	// var fractal lib.Fractal = lib.NewMandelbrotFractal(width, height, -0.7, 0.0, 4.0, 50)
	// Seahorse Valley:
	// var fractal = lib.NewMandelbrotFractal(width, height, -0.87591, 0.20464, 0.53184, 100)
	// Fractal from preset:
	for i, p := range presets.ColorPresets {
		fractPreset := presets.FractalPresets[8]
		width := 1920
		height := 1200
		fractal, err := lib.NewFractalFromPresets(width, height, p, fractPreset)
		if err != nil {
			panic(err)
		}
		img := lib.CalcFractalImage(fractal)
		img.SavePng(fmt.Sprintf("demo_images/image-color-preset-%02d.png", i))
	}
}
