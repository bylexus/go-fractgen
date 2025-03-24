package lib

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

type ColorPreset struct {
	Name    string
	Palette ColorPalette `json:"colors"`
}

type FractalPreset struct {
	Name               string
	PresetFunctionName string `json:"iterFunc"`
	ImageWidth         int    `json:"picWidth"`
	ImageHeight        int    `json:"picHeight"`
	DiameterCX         float64
	CenterCX           float64
	CenterCY           float64
	ColorPreset        string
	JuliaKi            float64
	JuliaKr            float64
	MaxIterations      int
	ColorPaletteLength int
	ColorPaletteRepeat int
}

func (f FractalPreset) FractalFunction() (FractalType, error) {

	switch f.PresetFunctionName {
	case "Mandelbrot":
		return FRACTAL_TYPE_MANDELBROT, nil
	case "Mandelbrot Z^3":
		return FRACTAL_TYPE_MANDELBROT3, nil
	case "Mandelbrot Z^4":
		return FRACTAL_TYPE_MANDELBROT4, nil
	case "Julia":
		return FRACTAL_TYPE_JULIA, nil
	default:
		return -1, errors.New("unknown fractal function")
	}
}

func ReadPresetJson(path string) ([]ColorPreset, []FractalPreset) {
	var colorPresets []ColorPreset = make([]ColorPreset, 0)
	var fractalPresets []FractalPreset = make([]FractalPreset, 0)

	jsonFile, err := os.Open(path)
	if err != nil {
		return colorPresets, fractalPresets
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
		return colorPresets, fractalPresets
	}

	var presets struct {
		ColorPresets   []ColorPreset   `json:"colorPresets"`
		FractalPresets []FractalPreset `json:"fractalPresets"`
	}

	err = json.Unmarshal(jsonData, &presets)
	if err != nil {
		log.Println(err)
		return colorPresets, fractalPresets
	}

	colorPresets = presets.ColorPresets
	fractalPresets = presets.FractalPresets

	return colorPresets, fractalPresets
}
