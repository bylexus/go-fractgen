package lib

import (
	"math"

	"github.com/bylexus/go-stdlib/ethreads"
)

type MandelbrotFractal struct {
	CommonFractParams
}

func NewMandelbrotFractal(fractalParams CommonFractParams) MandelbrotFractal {
	var params = initializeFractParams(fractalParams)
	return MandelbrotFractal{params}
}

func (f MandelbrotFractal) CreatePixelCalcJobFn(startPixX, startPixY, width, height int, img *FractImage) ethreads.JobFn {
	return func(id ethreads.ThreadId) {
		for pixY := startPixY; pixY < startPixY+height; pixY++ {
			for pixX := startPixX; pixX < startPixX+width; pixX++ {
				cx, cy := f.PixelToFractal(pixX, pixY)
				var fractRes FractFunctionResult
				fractRes = Mandelbrot(cx, cy, f.MaxAbsSquareAmount, f.MaxIterations)
				setImagePixel(img, pixX, pixY, f.CommonFractParams, fractRes)
			}
		}
	}
}

func (f MandelbrotFractal) ImageWidth() int {
	return f.CommonFractParams.ImageWidth
}

func (f MandelbrotFractal) ImageHeight() int {
	return f.CommonFractParams.ImageHeight
}

/*
An implementing algorithm for the fractal function: Mandelbrot set.

The Mandelbrot set is defined by:

	Z(n+1) = Z(n)^2 + c

while c = a constant complex number (cx + (cy)i) and Z(0) = 0;

cx = initial real value, calculated from the actual pixel's x position
cy = initial imaginary value, calculated from the actual pixel's y position

The number is iterated as long as it is clear that is is either reaching the border |Z^2| > max
or the max. number of iterations is reached.

@author Alexander Schenkel, www.alexi.ch
(c) 2012-2025 Alexander Schenkel
*/
func Mandelbrot(cx, cy, max_betrag_quadrat float64, maxIter int) FractFunctionResult {
	var betragQuadrat float64 = 0.0
	var iter int = 0
	var x, xt float64 = 0.0, 0.0
	var y, yt float64 = 0.0, 0.0
	var minDist = math.MaxFloat64
	var dist float64 = 0.0

	for betragQuadrat <= max_betrag_quadrat && iter < maxIter {
		xt = x*x - y*y + cx
		yt = 2*x*y + cy

		x = xt
		y = yt
		iter += 1
		betragQuadrat = x*x + y*y
		// point distance to point:
		// dist = math.Sqrt((x-2)*(x-2) + y*y)

		// dist = math.Abs((y - 2*x - 1) / math.Sqrt(2*2+1))
		// line equation ax + b:
		// point distance to line:
		//  dist = abs((a*x - y + b) / math.Sqrt(a*a + 1))
		// where x, y is the point and a, b are the coefficients of the line equation ax + y = 0

		dist = math.Abs((0.5*x - y + -3) / math.Sqrt(0.5*0.5+1))

		// dist = x
		if dist < minDist {
			minDist = dist
		}
	}
	result := FractFunctionResult{
		Iterations:   iter,
		BailoutValue: betragQuadrat,
		// for orbit trap:
		// BailoutValue: minDist,
		// BailoutValue: math.Sqrt(minDist),
	}
	return result
}
