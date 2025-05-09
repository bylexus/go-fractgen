package lib

import (
	"math/big"

	"github.com/bylexus/go-stdlib/ethreads"
)

type MandelbrotFractal struct {
	CommonFractParams
}

func NewMandelbrotFractal(fractalParams CommonFractParams) MandelbrotFractal {
	var params = initializeFractParams(fractalParams)
	return MandelbrotFractal{params}
}

func (f MandelbrotFractal) CreatePixelCalcFunc(pixX, pixY int, img *FractImage) ethreads.JobFn {
	return func(id ethreads.ThreadId) {
		cx, cy := f.PixelToFractal(pixX, pixY)
		var fractRes FractFunctionResult
		fractRes = Mandelbrot(cx, cy, BIG_MAX_ABS_SQUARE_AMOUNT, f.MaxIterations)
		setImagePixel(img, pixX, pixY, f.CommonFractParams, fractRes)
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
func Mandelbrot(cx, cy, max_betrag_quadrat *big.Float, maxIter int) FractFunctionResult {
	var betragQuadrat *big.Float = new(big.Float).SetPrec(SYS_PRECISION)
	var iter int = 0
	var x *big.Float = new(big.Float).SetPrec(SYS_PRECISION)
	var y *big.Float = new(big.Float).SetPrec(SYS_PRECISION)

	// for betragQuadrat <= max_betrag_quadrat && iter < maxIter {
	for betragQuadrat.Cmp(max_betrag_quadrat) <= 0 && iter < maxIter {
		// x*x
		xx := new(big.Float).Copy(x)
		xx.Mul(xx, x)

		// y*y
		yy := new(big.Float).Copy(y)
		yy.Mul(yy, y)

		// 2*x*y
		xy2 := new(big.Float).Copy(x)
		xy2.Mul(xy2, y)
		xy2.Mul(xy2, BIG_2)

		// xt = x*x - y*y + cx
		xt := new(big.Float).Copy(xx)
		xt.Sub(xt, yy).Add(xt, cx)

		// yt = 2*x*y + cy
		yt := xy2.Add(xy2, cy)

		x = xt
		y = yt
		iter += 1
		betragQuadrat.Mul(xx, yy)
	}
	result := FractFunctionResult{
		Iterations:   iter,
		BailoutValue: betragQuadrat,
	}
	return result
}
