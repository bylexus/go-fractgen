package lib

import (
	"image/color"
	"math"
)

type ColorPalette []color.RGBA

func setImagePixel(img *FractImage, x, y int, fractParams CommonFractParams, fractRes FractFunctionResult) {
	var LOG_2 float64 = math.Log(2)
	var LOG_MAX_BETRAG = math.Log(MAX_BETRAG_QUADRAT)
	// var LOG_4 float64 = math.Log(4)

	var iterValue float64
	// var whiteIntensity uint8
	if fractRes.Iterations <= fractParams.MaxIterations {
		if fractParams.SmoothColors == true {
			// Smooth coloring, see http://de.wikipedia.org/wiki/Mandelbrot-Menge#Iteration_eines_Bildpunktes:
			iterValue = float64(fractRes.Iterations) - math.Log(math.Log(fractRes.BailoutValue)/LOG_MAX_BETRAG)/LOG_2
		} else {
			// Rough coloring: Escape time algorithm:
			iterValue = float64(fractRes.Iterations)
		}
		// whiteIntensity = uint8(iterValue * 255 / float64(fractParams.MaxIterations))
		// whiteIntensity = uint8(iterValue * 255 / float64(fractParams.MaxIterations))

		// } else {
		// whiteIntensity = 0
	}
	// var whiteIntensity = uint8(math.Mod(iterValue, 255))
	// img.Set(x, y, color.RGBA{whiteIntensity, whiteIntensity, whiteIntensity, 255})

	setPaletteColor(img, x, y, iterValue, fractParams, fractRes)
}

func setPaletteColor(img *FractImage, x, y int, iterValue float64, fractParams CommonFractParams, fractRes FractFunctionResult) {
	// the color palette is defined with a set of "anchor colors": The full color palette is the same length
	// as the Max iteration count, but we only define anchor colors in-between. We then linear-interpolate the
	// correct color based on the iteration count-to-max-iterations ratio.

	var palette = fractParams.ColorPalette

	// Default: If the iteration count exceeds the max iterations, set the pixel color to black
	var selectedColor color.Color = color.RGBA{R: 0, G: 0, B: 0, A: 255}

	if iterValue <= float64(fractParams.MaxIterations) {
		// Calculate the iteration count-to-max-iterations ratio

		// Find the two anchor colors that the current iteration count is between
		if iterValue == float64(fractParams.MaxIterations) {
			selectedColor = palette[len(palette)-1]
		} else {
			ratio := float64(iterValue) / float64(fractParams.MaxIterations)
			colorPaletteSectionWidth := float64(fractParams.MaxIterations) / float64(len(palette)-1)

			var lower, upper color.RGBA
			var chosenPaletteIndex int = 0
			for i := 0; i < len(palette)-1; i++ {
				if float64(i)/float64(len(palette)-1) <= ratio {
					lower = palette[i]
					upper = palette[i+1]
					chosenPaletteIndex = i
				}
			}
			relativeColorDistance := iterValue - float64(chosenPaletteIndex)*colorPaletteSectionWidth
			relativeRatio := relativeColorDistance / colorPaletteSectionWidth

			// Linear-interpolate the correct color based on the ratio
			selectedColor = color.RGBA{
				R: uint8(float64(lower.R) + (float64(upper.R)-float64(lower.R))*relativeRatio),
				G: uint8(float64(lower.G) + (float64(upper.G)-float64(lower.G))*relativeRatio),
				B: uint8(float64(lower.B) + (float64(upper.B)-float64(lower.B))*relativeRatio),
				A: 255,
			}
		}

	}
	// Set the pixel color
	img.Set(x, y, selectedColor)
}
