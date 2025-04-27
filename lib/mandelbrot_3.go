package lib

import (
	"github.com/bylexus/go-stdlib/ethreads"
)

type Mandelbrot3Fractal struct {
	CommonFractParams
}

func NewMandelbrot3Fractal(fractalParams CommonFractParams) Mandelbrot3Fractal {
	var params = initializeFractParams(fractalParams)
	return Mandelbrot3Fractal{params}
}

func (f Mandelbrot3Fractal) CreatePixelCalcFunc(pixX, pixY int, img *FractImage) ethreads.JobFn {
	return func(id ethreads.ThreadId) {
		cx, cy := f.PixelToFractal(pixX, pixY)
		fractRes := Mandelbrot3(cx, cy, f.MaxAbsSquareAmount, f.MaxIterations)
		setImagePixel(img, pixX, pixY, f.CommonFractParams, fractRes)
	}
}

func (f Mandelbrot3Fractal) ImageWidth() int {
	return f.CommonFractParams.ImageWidth
}

func (f Mandelbrot3Fractal) ImageHeight() int {
	return f.CommonFractParams.ImageHeight
}

/*
An implementing algorithm for the fractal function: Mandelbrot set.

The Mandelbrot set is defined by:

	Z(n+1) = Z(n)^3 + c

while c = a constant complex number (cx + (cy)i) and Z(0) = 0;

cx = initial real value, calculated from the actual pixel's x position
cy = initial imaginary value, calculated from the actual pixel's y position

The number is iterated as long as it is clear that is is either reaching the border |Z^2| > max
or the max. number of iterations is reached.

@author Alexander Schenkel, www.alexi.ch
(c) 2012-2025 Alexander Schenkel
*/
func Mandelbrot3(cx, cy, max_betrag_quadrat float64, maxIter int) FractFunctionResult {
	var betragQuadrat float64 = 0.0
	var iter int = 0
	var x, xt float64 = 0.0, 0.0
	var y, yt float64 = 0.0, 0.0

	for betragQuadrat <= max_betrag_quadrat && iter < maxIter {
		// Z^3 + c:
		xt = x*(x*x-3*y*y) + cx
		yt = y*(3*x*x-y*y) + cy

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
