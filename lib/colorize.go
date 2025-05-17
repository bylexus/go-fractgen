package lib

import (
	"image/color"
	"math"
)

type PaletteEntry struct {
	color.RGBA
	Steps int `json:"steps"`
}

type ColorPalette []PaletteEntry

const defaultPaletteLength = 256

func setImagePixel(img *FractImage, x, y int, fractParams CommonFractParams, fractRes FractFunctionResult) {
	// calc the color for this pixel:
	var LOG_2 float64 = math.Log(2)
	var LOG_MAX_BETRAG = math.Log(fractParams.MaxAbsSquareAmount)
	// var LOG_4 float64 = math.Log(4)

	var iterValue float64
	if fractRes.Iterations <= fractParams.MaxIterations {
		if fractParams.SmoothColors == true {
			// Smooth coloring, see http://de.wikipedia.org/wiki/Mandelbrot-Menge#Iteration_eines_Bildpunktes:
			iterValue = float64(fractRes.Iterations) - math.Log(math.Log(fractRes.BailoutValue)/LOG_MAX_BETRAG)/LOG_2
		} else {
			// Rough coloring: Escape time algorithm:
			iterValue = float64(fractRes.Iterations)
		}
	}

	SetPaletteColor(img, x, y, iterValue, fractParams, fractRes)
}

func SetPaletteColor(img *FractImage, x, y int, iterValue float64, fractParams CommonFractParams, fractRes FractFunctionResult) {
	// the color palette is defined with a set of "anchor colors": The full color palette is the same length
	// as the Max iteration count, but we only define anchor colors in-between. We then linear-interpolate the
	// correct color based on the iteration count-to-max-iterations ratio.

	var palette = fractParams.ColorPalette
	var nrOfColors = len(palette)
	if nrOfColors == 0 {
		return
	}

	// Default: If the iteration count exceeds the max iterations, set the pixel color to black
	var selectedColor PaletteEntry = PaletteEntry{color.RGBA{R: 0, G: 0, B: 0, A: 255}, 255}
	var repeat = 1
	var maxIterations = float64(fractParams.MaxIterations)
	if fractParams.ColorPaletteRepeat > 0 {
		repeat = fractParams.ColorPaletteRepeat
	}
	/*
			A palette is defined from single color stops. Each stop also defines its "length", ranging from 0 (=0%) to 256 (=100%):

			|----------|-----|----------|-----|
			A:256      B:128 C:256      D:128
			| 100%     | 50% | 100%     |  50%|

			The color corresponds to the iter value, from 0 to Max-Iter:
			|----------|-----|----------|-----|
		 .  0                                 Max-Iter

		    We calculate the "palette length" by summing up the color stops' lengths, which accounts for 100% of the length.
			Then we create an array with the color stop's length, to determine the correct color for each pixel, interpolating
			the colors between 2 color stops (linear interpolation).

			A palette repeat simply repeats the palette n times, so that the palette fits multiple times within
			the max iter boundary.

			Example: We have the following palette:
			[{col1, steps: 256}, {col2, steps: 128}, {col3, steps: 256}, {col4, steps: 128}]
			This sums up to a total length of 768 (2*256 + 2*128), with a corresponding
			color stop array:
			[256, 128, 256, 128]

			Then we have the following iteration values:
			- Iterations for the pixel: 880
			- Max. Iterations: 1000

			This means we need the color stop 768 * (800/1000) = 675 (total palette length * percentage of max iterations)
			So we can sum up from our color stop array until we reach the color stop larger than 675, and we then
			know the lower and upper index in the color stop array, that corresponds to the 2 palette entries we need.

			In our example:
			The palette entry for index 675 leads to the 4th color stop entry (last 128, as 256+128+258+128 > 675): That means
			that our color is between the 4th palette entry and the 5th (=0th, as we wrap around).
			The percentage within this two palettes then is calculated and interpolated:

			(675 - 256 - 128 - 256) / 128 = 0.27 --> 27% from the left color stop.

	*/

	if iterValue <= maxIterations {
		if fractParams.ColorPaletteLength > 0 {
			// means: the palette length is defined by the user, not by the max iterations
			maxIterations = float64(fractParams.ColorPaletteLength)
			iterValue = math.Mod(iterValue, maxIterations)
			repeat = 1
		}

		// reverse value to take the color from the other side: same as reversing the palette:
		if fractParams.ColorPaletteReverse {
			iterValue = maxIterations - iterValue
		}

		// var paletteLength = len(palette) * repeat
		var paletteLength = 0
		var colorStopLengths = make([]int, 0)
		for i := 0; i < repeat; i++ {
			for _, entry := range palette {
				if entry.Steps < 0 {
					paletteLength += 0
					colorStopLengths = append(colorStopLengths, 0)
				} else if entry.Steps > 0 {
					paletteLength += entry.Steps
					colorStopLengths = append(colorStopLengths, entry.Steps)
				} else {
					paletteLength += defaultPaletteLength
					colorStopLengths = append(colorStopLengths, defaultPaletteLength)
				}
			}
		}

		// Find the two anchor colors that the current iteration count is between
		ratio := iterValue / float64(maxIterations)
		paletteEntry := ratio * float64(paletteLength) // Position within the full palette

		// colorPaletteSectionWidth := float64(maxIterations) / float64(paletteLength)

		var lower, upper PaletteEntry
		// var chosenPaletteIndex int = 0
		var stopsUntilNow = 0
		var upperStop = 0
		var actStopWidth = 0
		lower = palette[0]
		upper = palette[1%nrOfColors]
		for i, entry := range colorStopLengths {
			upperStop += entry
			actStopWidth = entry
			if paletteEntry > float64(upperStop) {
				stopsUntilNow = upperStop
				continue
			} else {
				lower = palette[i%nrOfColors]
				upper = palette[(i+1)%nrOfColors]
				break
			}
		}

		if fractParams.ColorPaletteHardStops {
			selectedColor = lower
		} else {
			// Linear-interpolate the correct color based on the ratio
			relativeRatio := (paletteEntry - float64(stopsUntilNow)) / float64(actStopWidth)
			selectedColor = PaletteEntry{color.RGBA{
				R: uint8(math.Round(float64(lower.R) + (float64(upper.R)-float64(lower.R))*relativeRatio)),
				G: uint8(math.Round(float64(lower.G) + (float64(upper.G)-float64(lower.G))*relativeRatio)),
				B: uint8(math.Round(float64(lower.B) + (float64(upper.B)-float64(lower.B))*relativeRatio)),
				A: 255,
			}, 255}
		}
	}
	// Set the pixel color
	img.Set(x, y, selectedColor.RGBA)
}
