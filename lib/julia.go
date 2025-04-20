package lib

import (
	"github.com/bylexus/go-stdlib/ethreads"
)

type JuliaFractal struct {
	CommonFractParams
	JuliaKr float64
	JuliaKi float64
}

func NewJuliaFractal(imageWidth, imageHeight int, centerCX, centerCY, diameterCX float64, maxIterations int, colorPalette ColorPalette, colorPaletteRepeat int, juliaKr, juliaKi float64) JuliaFractal {
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

	return JuliaFractal{params, juliaKr, juliaKi}
}

func (f JuliaFractal) ImageWidth() int {
	return f.CommonFractParams.ImageWidth
}

func (f JuliaFractal) ImageHeight() int {
	return f.CommonFractParams.ImageHeight
}

func (f JuliaFractal) CreatePixelCalcFunc(pixX, pixY int, img *FractImage) ethreads.JobFn {
	return func(id ethreads.ThreadId) {
		cx, cy := f.PixelToFractal(pixX, pixY)
		fractRes := Julia(cx, cy, f.MaxAbsSquareAmount, f.MaxIterations, f.JuliaKr, f.JuliaKi)
		setImagePixel(img, pixX, pixY, f.CommonFractParams, fractRes)
	}
}

/*
*
  - An implementing algorithm for the fractal function: Julia set.
    *
  - The julia set is defined by:
    *
  - Z(n+1) = Z(n)^2 + K
    *
  - while K = a constant complex number (e.g. -0.6 + 0.6i)
  - and
  - Z(0) = (cx + (cy)i) + K;
  - cx = initial real value, calculated from the actual pixel's x position
  - cy = initial imaginary value, calculated from the actual pixel's y position
    *
  - The number is iterated as long as it is clear that is is either reaching the border |Z^2| > max
  - or the max. number of iterations is reached.
    *
  - Part of JFractGen - a Julia / Mandelbrot Fractal generator written in Java/Swing.
  - @author Alexander Schenkel, www.alexi.ch
  - (c) 2012 Alexander Schenkel
*/
func Julia(cx, cy, max_betrag_quadrat float64, maxIter int, julia_r, julia_i float64) FractFunctionResult {
	var betragQuadrat float64 = 0.0
	var iter int = 0
	var x, xt float64 = cx, 0.0
	var y, yt float64 = cy, 0.0

	for betragQuadrat <= max_betrag_quadrat && iter < maxIter {
		xt = x*x - y*y + julia_r
		yt = 2*x*y + julia_i

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
