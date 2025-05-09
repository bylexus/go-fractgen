package lib

import (
	"math/big"

	"github.com/bylexus/go-stdlib/ethreads"
)

type JuliaFractal struct {
	CommonFractParams
	JuliaKr *big.Float
	JuliaKi *big.Float
}

func NewJuliaFractal(fractalParams CommonFractParams, juliaKr, juliaKi *big.Float) JuliaFractal {
	var params = initializeFractParams(fractalParams)

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
		fractRes := Julia(cx, cy, BIG_MAX_ABS_SQUARE_AMOUNT, f.MaxIterations, f.JuliaKr, f.JuliaKi)
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
func Julia(cx, cy, max_betrag_quadrat *big.Float, maxIter int, julia_r, julia_i *big.Float) FractFunctionResult {
	var betragQuadrat *big.Float = new(big.Float).SetPrec(SYS_PRECISION)
	var iter int = 0
	var x *big.Float = new(big.Float).SetPrec(SYS_PRECISION)
	var y *big.Float = new(big.Float).SetPrec(SYS_PRECISION)

	// start value of x/y:
	x.Copy(cx)
	y.Copy(cy)

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

		// xt = x*x - y*y + julia_r
		xt := new(big.Float).Copy(xx)
		xt.Sub(xt, yy)
		xt.Add(xt, julia_r)

		// yt = 2*x*y + julia_i
		yt := new(big.Float).Copy(xy2)
		yt.Add(yt, julia_i)

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
