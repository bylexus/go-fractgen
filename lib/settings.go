package lib

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"strings"
)

type ColorPresets []ColorPreset
type FractalPresets []FractalPreset

type Presets struct {
	ColorPresets   ColorPresets   `json:"colorPresets"`
	FractalPresets FractalPresets `json:"fractalPresets"`
}

func (p ColorPresets) GetByIdent(ident string) (ColorPreset, error) {
	for _, preset := range p {
		if strings.ToLower(preset.Ident) == ident {
			return preset, nil
		}
	}
	return ColorPreset{}, errors.New("no color preset found")
}

func (p FractalPresets) GetByName(name string) (FractalPreset, error) {
	name = strings.ToLower(name)
	for _, preset := range p {
		if strings.ToLower(preset.Name) == name {
			return preset, nil
		}
	}
	return FractalPreset{}, errors.New("no fractal preset found")
}

type ColorPreset struct {
	Name    string       `json:"name"`
	Ident   string       `json:"ident"`
	Palette ColorPalette `json:"colors"`
}

type FractalPreset struct {
	Name                string  `json:"name"`
	IterFunc            string  `json:"iterFunc"`
	DiameterCX          float64 `json:"diameterCX"`
	CenterCX            float64 `json:"centerCX"`
	CenterCY            float64 `json:"centerCY"`
	ColorPreset         string  `json:"colorPreset"`
	JuliaKi             float64 `json:"juliaKi"`
	JuliaKr             float64 `json:"juliaKr"`
	MaxIterations       int     `json:"maxIterations"`
	ColorPaletteLength  int     `json:"colorPaletteLength"`
	ColorPaletteRepeat  int     `json:"colorPaletteRepeat"`
	ColorPaletteReverse bool    `json:"colorPaletteReverse"`
}

func (f FractalPreset) FractalFunction() (FractalType, error) {

	switch f.IterFunc {
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

func ReadPresetJson(path string) Presets {
	var presets Presets = Presets{
		ColorPresets:   make([]ColorPreset, 0),
		FractalPresets: make([]FractalPreset, 0),
	}

	jsonFile, err := os.Open(path)
	if err != nil {
		return presets
	}
	defer jsonFile.Close()

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
		return presets
	}

	err = json.Unmarshal(jsonData, &presets)
	if err != nil {
		log.Println(err)
		return presets
	}

	return presets
}
