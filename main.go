package main

import (
	"fmt"
	"log"

	"github.com/bylexus/go-fract/lib"
	"github.com/bylexus/go-fract/web"
)

func main() {

	// demo()
	webserver()

}

func webserver() {
	presetsFile := "presets.json"
	colorPresets, fractalPresets := lib.ReadPresetJson(presetsFile)

	var server *web.WebServer = web.NewWebServer(colorPresets, fractalPresets)

	fmt.Printf("Starting Webserver, listen on %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func demo() {
	presetsFile := "presets.json"

	colorPresets, fractalPresets := lib.ReadPresetJson(presetsFile)

	// Standard mandelbrot view:
	// var fractal lib.Fractal = lib.NewMandelbrotFractal(width, height, -0.7, 0.0, 4.0, 50)
	// Seahorse Valley:
	// var fractal = lib.NewMandelbrotFractal(width, height, -0.87591, 0.20464, 0.53184, 100)
	// Fractal from preset:
	for i, p := range colorPresets {
		fractPreset := fractalPresets[8]
		fractPreset.ImageWidth = 1920
		fractPreset.ImageHeight = 1280
		fractal, err := lib.NewFractalFromPresets(p, fractPreset)
		if err != nil {
			panic(err)
		}
		img := fractal.CalcFractalImage(nil)
		img.SavePng(fmt.Sprintf("demo_images/image-color-preset-%02d.png", i))
	}
}
