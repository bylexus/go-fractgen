package lib

import (
	"github.com/bylexus/go-stdlib/ethreads"
)

type Mandelbrot4Fractal struct {
	CommonFractParams
}

func NewMandelbrot4Fractal(imageWidth, imageHeight int, centerCX, centerCY, diameterCX float64, maxIterations int, colorPalette ColorPalette, colorPaletteRepeat int) Mandelbrot4Fractal {
	var params = initializeFractParams(CommonFractParams{
		ImageWidth:         imageWidth,
		ImageHeight:        imageHeight,
		CenterCX:           centerCX,
		CenterCY:           centerCY,
		DiameterCX:         diameterCX,
		MaxIterations:      maxIterations,
		ColorPalette:       colorPalette,
		ColorPaletteRepeat: colorPaletteRepeat,
	})

	return Mandelbrot4Fractal{params}
}

func (f Mandelbrot4Fractal) CreatePixelCalcFunc(pixX, pixY int, img *FractImage) ethreads.JobFn {
	return func(id ethreads.ThreadId) {
		cx, cy := f.PixelToFractal(pixX, pixY)
		fractRes := Mandelbrot4(cx, cy, f.MaxAbsSquareAmount, f.MaxIterations)
		setImagePixel(img, pixX, pixY, f.CommonFractParams, fractRes)
	}
}

func (f Mandelbrot4Fractal) ImageWidth() int {
	return f.CommonFractParams.ImageWidth
}

func (f Mandelbrot4Fractal) ImageHeight() int {
	return f.CommonFractParams.ImageHeight
}

/*
An implementing algorithm for the fractal function: Mandelbrot set.

The Mandelbrot set is defined by:

	Z(n+1) = Z(n)^4 + c

while c = a constant complex number (cx + (cy)i) and Z(0) = 0;

cx = initial real value, calculated from the actual pixel's x position
cy = initial imaginary value, calculated from the actual pixel's y position

The number is iterated as long as it is clear that is is either reaching the border |Z^2| > max
or the max. number of iterations is reached.

@author Alexander Schenkel, www.alexi.ch
(c) 2012-2025 Alexander Schenkel
*/
func Mandelbrot4(cx, cy, max_betrag_quadrat float64, maxIter int) FractFunctionResult {
	var betragQuadrat float64 = 0.0
	var iter int = 0
	var x, xt float64 = 0.0, 0.0
	var y, yt float64 = 0.0, 0.0

	for betragQuadrat <= max_betrag_quadrat && iter < maxIter {
		// Z^4 + c:
		xt = x*x*x*x - 6*x*x*y*y + y*y*y*y + cx
		yt = 4*x*x*x*y - 4*x*y*y*y + cy

		x = xt
		y = yt
		iter += 1
		betragQuadrat = x*x + y*y
	}
	result := FractFunctionResult{
		Iterations:   iter,
		BailoutValue: betragQuadrat,
	}
	return result
}
