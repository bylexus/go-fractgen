package lib

import (
	"errors"
	"runtime"

	"github.com/bylexus/go-stdlib/ethreads"
)

// const MAX_BETRAG_QUADRAT float64 = 256

const MAX_BETRAG_QUADRAT float64 = 4

type FractalType int

const (
	FRACTAL_TYPE_MANDELBROT = iota
	FRACTAL_TYPE_MANDELBROT2
	FRACTAL_TYPE_MANDELBROT3
	FRACTAL_TYPE_MANDELBROT4
	FRACTAL_TYPE_JULIA
)

type Fractal interface {
	CalcFractalImage(threadPool *ethreads.ThreadPool) *FractImage
}

type FractFunctionResult struct {
	Iterations   int
	BailoutValue float64
}

type CommonFractParams struct {
	MaxIterations int

	CenterCX   float64
	CenterCY   float64
	DiameterCX float64

	// JuliaKr float64
	// JuliaKi float64

	ImageWidth  int
	ImageHeight int

	SmoothColors       bool
	FixedSizePalette   bool
	ColorPaletteRepeat int

	// calculaed during initialization:
	aspect float64
	minCX  float64
	maxCX  float64
	minCY  float64
	maxCY  float64

	ColorPalette ColorPalette
}

func (f CommonFractParams) PixelToFractal(x, y int) (cx, cy float64) {
	// y axis is inverted in image and fractal space:
	y = f.ImageHeight - y
	cx = f.minCX + (f.maxCX-f.minCX)*(float64(x)/float64(f.ImageWidth))
	cy = f.minCY + (f.maxCY-f.minCY)*(float64(y)/float64(f.ImageHeight))
	return cx, cy
}

func NewMandelbrotFractal(imageWidth, imageHeight int, centerCX, centerCY, diameterCX float64, maxIterations int, colorPalette ColorPalette, colorPaletteRepeat int) MandelbrotFractal {
	var aspect, fract_width, fract_heigth float64

	aspect = float64(imageWidth) / float64(imageHeight)
	fract_width = diameterCX
	fract_heigth = diameterCX / aspect

	var min_cx = centerCX - (fract_width / 2)
	var max_cx = min_cx + fract_width
	var min_cy = centerCY - (fract_heigth / 2)
	var max_cy = min_cy + fract_heigth

	if colorPaletteRepeat <= 0 {
		colorPaletteRepeat = 1
	}

	var params CommonFractParams = CommonFractParams{
		MaxIterations: maxIterations,
		CenterCX:      centerCX,
		CenterCY:      centerCY,
		DiameterCX:    diameterCX,
		// JuliaKr:          -0.6,
		// JuliaKi:          0.6,
		ImageWidth:         imageWidth,
		ImageHeight:        imageHeight,
		SmoothColors:       true,
		FixedSizePalette:   false,
		ColorPaletteRepeat: colorPaletteRepeat,

		ColorPalette: colorPalette,

		// Calculated during initialization:
		aspect: aspect,
		minCX:  min_cx,
		maxCX:  max_cx,
		minCY:  min_cy,
		maxCY:  max_cy,
	}

	return MandelbrotFractal{params}
}

type MandelbrotFractal struct {
	CommonFractParams
}

func (f MandelbrotFractal) CalcFractalImage(threadPool *ethreads.ThreadPool) *FractImage {
	if threadPool == nil {
		threadPool := ethreads.NewThreadPool(runtime.NumCPU()*2, nil)
		threadPool.Start()
		defer threadPool.Shutdown()
	}

	img := NewFractImage(f.ImageWidth, f.ImageHeight)

	for y := 0; y < f.ImageHeight; y++ {
		for x := 0; x < f.ImageWidth; x++ {
			f := func(pixX, pixY int) ethreads.JobFn {
				return func(id ethreads.ThreadId) {
					cx, cy := f.PixelToFractal(pixX, pixY)
					fractRes := Mandelbrot(cx, cy, MAX_BETRAG_QUADRAT, f.MaxIterations)
					setImagePixel(img, pixX, pixY, f.CommonFractParams, fractRes)
				}
			}
			threadPool.AddJobFn(f(x, y))
		}
	}

	return img
}

func NewFractalFromPresets(colorPreset ColorPreset, fractalPreset FractalPreset) (Fractal, error) {
	fractFunc, err := fractalPreset.FractalFunction()
	if err != nil {
		return nil, err
	}
	switch fractFunc {
	case FRACTAL_TYPE_MANDELBROT:
		return NewMandelbrotFractal(
			fractalPreset.ImageWidth,
			fractalPreset.ImageHeight,
			fractalPreset.CenterCX,
			fractalPreset.CenterCY,
			fractalPreset.DiameterCX,
			fractalPreset.MaxIterations,
			colorPreset.Palette,
			fractalPreset.ColorPaletteRepeat,
		), nil

	default:
		return nil, errors.New("unknown fractal function")
	}
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

Part of JFractGen - a Julia / Mandelbrot Fractal generator written in Java/Swing.
@author Alexander Schenkel, www.alexi.ch
(c) 2012 Alexander Schenkel
*/
func Mandelbrot(cx, cy, max_betrag_quadrat float64, maxIter int) FractFunctionResult {
	var betragQuadrat float64 = 0.0
	var iter int = 0
	var x, xt float64 = 0.0, 0.0
	var y, yt float64 = 0.0, 0.0

	for betragQuadrat <= max_betrag_quadrat && iter < maxIter {
		xt = x*x - y*y + cx
		yt = 2*x*y + cy

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

/**

/**
 * An implementing algorithm for the fractal function: Mandelbrot set.
 *
 * The Mandelbrot set is defined by:
 *
 * Z(n+1) = Z(n)^2 + c
 *
 * while c = a constant complex number (cx + (cy)i)
 * and
 * Z(0) = 0;
 *
 * cx = initial real value, calculated from the actual pixel's x position
 * cy = initial imaginary value, calculated from the actual pixel's y position
 *
 * The number is iterated as long as it is clear that is is either reaching the border |Z^2| > max
 * or the max. number of iterations is reached.
 *
 * Part of JFractGen - a Julia / Mandelbrot Fractal generator written in Java/Swing.
 * @author Alexander Schenkel, www.alexi.ch
 * (c) 2012 Alexander Schenkel
	public FractFunctionResult fractIterFunc(double cx, double cy, double max_betrag_quadrat,
			double max_iter, double julia_r, double julia_i) {
		double betrag_quadrat = 0;
		double iter = 0;
		double x = 0, xt;
		double y = 0, yt;

		while (betrag_quadrat <= max_betrag_quadrat && iter < max_iter) {
			xt = x * x - y*y + cx;
			yt = 2*x*y + cy;

			// Z^3 + c:
			//xt = x*(x*x - 3*y*y) + cx;
			//yt = y*(3*x*x - y*y) + cy;

			// Z^4 + c:
			//xt = x*x*x*x -6*x*x*y*y + y*y*y*y + cx;
			//yt = 4*x*x*x*y - 4*x*y*y*y + cy;
			x = xt;
			y = yt;
			iter += 1;
			betrag_quadrat = x*x + y*y;
		}
		FractFunctionResult r = new FractFunctionResult();
		r.iterValue = iter;
		r.bailoutValue = betrag_quadrat;
		return r;
	}


	// --------------- Mandelbrot ^3 ----------------
	/**
 * An implementing algorithm for the fractal function: Mandelbrot set.
 *
 * This Mandelbrot set is defined by:
 *
 * Z(n+1) = Z(n)^3 + c
 *
 * while c = a constant complex number (cx + (cy)i)
 * and
 * Z(0) = 0;
 *
 * cx = initial real value, calculated from the actual pixel's x position
 * cy = initial imaginary value, calculated from the actual pixel's y position
 *
 * The number is iterated as long as it is clear that is is either reaching the border |Z^3| > max
 * or the max. number of iterations is reached.
 *
 * Part of JFractGen - a Julia / Mandelbrot Fractal generator written in Java/Swing.
 * @author Alexander Schenkel, www.alexi.ch
 * (c) 2012 Alexander Schenkel
	public FractFunctionResult fractIterFunc(double cx, double cy, double max_betrag_quadrat,
			double max_iter, double julia_r, double julia_i) {
		double betrag_quadrat = 0;
		double iter = 0;
		double x = 0, xt;
		double y = 0, yt;

		while (betrag_quadrat <= max_betrag_quadrat && iter < max_iter) {
			// Z^3 + c:
			xt = x*(x*x - 3*y*y) + cx;
			yt = y*(3*x*x - y*y) + cy;
			x = xt;
			y = yt;
			iter += 1;
			betrag_quadrat = x*x + y*y;
		}
		FractFunctionResult r = new FractFunctionResult();
		r.iterValue = iter;
		r.bailoutValue = betrag_quadrat;
		return r;
	}




------------ Mandelbrot ^4 ----------------
/**
 * An implementing algorithm for the fractal function: Mandelbrot set.
 *
 * This Mandelbrot set is defined by:
 *
 * Z(n+1) = Z(n)^4 + c
 *
 * while c = a constant complex number (cx + (cy)i)
 * and
 * Z(0) = 0;
 *
 * cx = initial real value, calculated from the actual pixel's x position
 * cy = initial imaginary value, calculated from the actual pixel's y position
 *
 * The number is iterated as long as it is clear that is is either reaching the border |Z^4| > max
 * or the max. number of iterations is reached.
 *
 * Part of JFractGen - a Julia / Mandelbrot Fractal generator written in Java/Swing.
 * @author Alexander Schenkel, www.alexi.ch
 * (c) 2012 Alexander Schenkel
 public class Mandelbrot4FractFunction implements IFractFunction {
	@Override
	public String toString() {
		return "Mandelbrot Z^4";
	}

	public FractFunctionResult fractIterFunc(double cx, double cy, double max_betrag_quadrat,
			double max_iter, double julia_r, double julia_i) {
		double betrag_quadrat = 0;
		double iter = 0;
		double x = 0, xt;
		double y = 0, yt;

		while (betrag_quadrat <= max_betrag_quadrat && iter < max_iter) {
			// Z^4 + c:
			xt = x*x*x*x -6*x*x*y*y + y*y*y*y + cx;
			yt = 4*x*x*x*y - 4*x*y*y*y + cy;
			x = xt;
			y = yt;
			iter += 1;
			betrag_quadrat = x*x + y*y;
		}
		FractFunctionResult r = new FractFunctionResult();
		r.iterValue = iter;
		r.bailoutValue = betrag_quadrat;
		return r;
	}


-------------- Julia ----------------
/**
 * An implementing algorithm for the fractal function: Julia set.
 *
 * The julia set is defined by:
 *
 * Z(n+1) = Z(n)^2 + K
 *
 * while K = a constant complex number (e.g. -0.6 + 0.6i)
 * and
 * Z(0) = (cx + (cy)i) + K;
 * cx = initial real value, calculated from the actual pixel's x position
 * cy = initial imaginary value, calculated from the actual pixel's y position
 *
 * The number is iterated as long as it is clear that is is either reaching the border |Z^2| > max
 * or the max. number of iterations is reached.
 *
 * Part of JFractGen - a Julia / Mandelbrot Fractal generator written in Java/Swing.
 * @author Alexander Schenkel, www.alexi.ch
 * (c) 2012 Alexander Schenkel
 public class JuliaFractFunction implements IFractFunction {
	public String toString() {
		return "Julia";
	}

	public FractFunctionResult fractIterFunc(double cx, double cy, double max_betrag_quadrat,
			double max_iter, double julia_r, double julia_i) {
		double betrag_quadrat = 0;
		double iter = 0;
		double x = cx,xt;
		double y = cy,yt;

		while (betrag_quadrat <= max_betrag_quadrat && iter < max_iter) {
			xt = x * x - y*y + julia_r;
			yt = 2*x*y + julia_i;
			x = xt;
			y = yt;
			iter += 1;
			betrag_quadrat = x*x + y*y;
		}
		FractFunctionResult r = new FractFunctionResult();
		r.iterValue = iter;
		r.bailoutValue = betrag_quadrat;
		return r;
	}
}

*/
