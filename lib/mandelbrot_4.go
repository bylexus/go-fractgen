package lib

import (
	"math/big"

	"github.com/bylexus/go-stdlib/ethreads"
)

type Mandelbrot4Fractal struct {
	CommonFractParams
}

func NewMandelbrot4Fractal(fractalParams CommonFractParams) Mandelbrot4Fractal {
	var params = initializeFractParams(fractalParams)

	return Mandelbrot4Fractal{params}
}

func (f Mandelbrot4Fractal) CreatePixelCalcFunc(pixX, pixY int, img *FractImage) ethreads.JobFn {
	return func(id ethreads.ThreadId) {
		cx, cy := f.PixelToFractal(pixX, pixY)
		fractRes := Mandelbrot4(cx, cy, BIG_MAX_ABS_SQUARE_AMOUNT, f.MaxIterations)
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
func Mandelbrot4(cx, cy, max_betrag_quadrat *big.Float, maxIter int) FractFunctionResult {
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

		// x^4
		xxxx := new(big.Float).Copy(xx)
		xxxx.Mul(xxxx, xx)

		// y^4
		yyyy := new(big.Float).Copy(yy)
		yyyy.Mul(yyyy, yy)

		// 6*x^2*y^2
		xxyy6 := new(big.Float).Copy(xx)
		xxyy6.Mul(xxyy6, yy)
		xxyy6.Mul(xxyy6, BIG_6)

		// 4*x^3*y
		xxxy4 := new(big.Float).Copy(xx)
		xxxy4.Mul(xxxy4, x)
		xxxy4.Mul(xxxy4, y)
		xxxy4.Mul(xxxy4, BIG_4)

		// 4*x*y^3
		xyyy4 := new(big.Float).Copy(x)
		xyyy4.Mul(xyyy4, yy)
		xyyy4.Mul(xyyy4, y)
		xyyy4.Mul(xyyy4, BIG_4)

		// Z^4 + c:
		// xt = x*x*x*x - 6*x*x*y*y + y*y*y*y + cx
		xt := xxxx.Sub(xxxx, xxyy6).Add(xxxx, yyyy).Add(xxxx, cx)

		// yt = 4*x*x*x*y - 4*x*y*y*y + cy
		yt := xxxy4.Sub(xxxy4, xyyy4).Add(xxxy4, cy)

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
