package lib

import (
	"errors"
	"math/big"
	"runtime"

	"github.com/bylexus/go-stdlib/ethreads"
)

// max amount = 256 gives smoother colors around the fractal:
// const MAX_ABS_SQUARE_AMOUNT float64 = 4
const MAX_ABS_SQUARE_AMOUNT float64 = 256

type FractalType string

const (
	FRACTAL_TYPE_MANDELBROT  = "mandelbrot"
	FRACTAL_TYPE_MANDELBROT3 = "mandelbrot3"
	FRACTAL_TYPE_MANDELBROT4 = "mandelbrot4"
	FRACTAL_TYPE_JULIA       = "julia"
)

type Fractal interface {
	CreatePixelCalcJobFn(startPixX, startPixY, width, height int, img *FractImage) ethreads.JobFn
	ImageWidth() int
	ImageHeight() int
}

type FractFunctionResult struct {
	Iterations   int
	BailoutValue float64
}

type CommonFractParams struct {
	MaxAbsSquareAmount float64
	MaxIterations      int

	CenterCX   float64
	CenterCY   float64
	DiameterCX float64

	ImageWidth  int
	ImageHeight int

	SmoothColors          bool
	ColorPaletteLength    int
	ColorPalette          ColorPalette
	ColorPaletteRepeat    int
	ColorPaletteReverse   bool
	ColorPaletteHardStops bool

	// calculaed during initialization:
	aspect float64
	minCX  float64
	maxCX  float64
	minCY  float64
	maxCY  float64
}

func (f CommonFractParams) PixelToFractal(x, y int) (cx, cy float64) {
	// y axis is inverted in image and fractal space:
	y = f.ImageHeight - y
	cx = f.minCX + (f.maxCX-f.minCX)*(float64(x)/float64(f.ImageWidth))
	cy = f.minCY + (f.maxCY-f.minCY)*(float64(y)/float64(f.ImageHeight))
	return cx, cy
}

func initializeFractParams(commonFractParams CommonFractParams) CommonFractParams {
	// var aspect, fract_width, fract_heigth float64
	var aspect, fract_width, fract_heigth, centerCX, centerCY *big.Float

	aspect = big.NewFloat(0).Quo(big.NewFloat(float64(commonFractParams.ImageWidth)), big.NewFloat(float64(commonFractParams.ImageHeight)))
	// aspect = float64(commonFractParams.ImageWidth) / float64(commonFractParams.ImageHeight)
	fract_width = big.NewFloat(float64(commonFractParams.DiameterCX))
	fract_heigth = big.NewFloat(0.0).Quo(fract_width, aspect)
	centerCX = big.NewFloat(commonFractParams.CenterCX)
	centerCY = big.NewFloat(commonFractParams.CenterCY)
	two := big.NewFloat(2.0)

	// var min_cx float64 = commonFractParams.CenterCX - (fract_width / 2.0)
	var min_cx = big.NewFloat(0.0).Sub(centerCX, (big.NewFloat(0.0).Quo(fract_width, two)))
	// var max_cx float64 = min_cx + fract_width
	var max_cx = big.NewFloat(0.0).Add(min_cx, fract_width)
	// var min_cy float64 = commonFractParams.CenterCY - (fract_heigth / 2.0)
	var min_cy = big.NewFloat(0.0).Sub(centerCY, (big.NewFloat(0.0).Quo(fract_heigth, two)))
	// var max_cy float64 = min_cy + fract_heigth
	var max_cy = big.NewFloat(0.0).Add(min_cy, fract_heigth)

	if commonFractParams.ColorPaletteRepeat <= 0 {
		commonFractParams.ColorPaletteRepeat = 1
	}
	if commonFractParams.ColorPaletteLength <= 0 {
		commonFractParams.ColorPaletteLength = -1
	}

	commonFractParams.SmoothColors = true

	// Calculated during initialization:
	commonFractParams.aspect, _ = aspect.Float64()
	commonFractParams.minCX, _ = min_cx.Float64()
	commonFractParams.maxCX, _ = max_cx.Float64()
	commonFractParams.minCY, _ = min_cy.Float64()
	commonFractParams.maxCY, _ = max_cy.Float64()
	commonFractParams.MaxAbsSquareAmount = MAX_ABS_SQUARE_AMOUNT

	return commonFractParams
}

func CalcFractalImage(f Fractal) *FractImage {
	tp := ethreads.NewThreadPool(runtime.NumCPU()*2, nil)
	tp.Start()

	img := NewFractImage(f.ImageWidth(), f.ImageHeight())

	// We calculate blocks of pixels in separate goroutines: For each block,
	// we start a new goroutine in the thread pool.
	// A single pixel per goroutine is too inperformant / generates too many goroutines.
	// A block size of 64x64 pixels is a good compromise between inperformant and too many goroutines.
	var blockWidth, blockHeight = 64, 64
	for y := 0; y < f.ImageHeight(); y += blockHeight {
		for x := 0; x < f.ImageWidth(); x += blockWidth {
			tp.AddJobFn(f.CreatePixelCalcJobFn(x, y, blockWidth, blockHeight, img))
		}
	}
	tp.Shutdown()

	return img
}

func NewFractalFromPresets(width, height int, colorPreset ColorPreset, fractalPreset FractalPreset) (Fractal, error) {
	fractFunc, err := fractalPreset.FractalFunction()
	if err != nil {
		return nil, err
	}
	commonParams := CommonFractParams{
		ImageWidth:            width,
		ImageHeight:           height,
		CenterCX:              fractalPreset.CenterCX,
		CenterCY:              fractalPreset.CenterCY,
		DiameterCX:            fractalPreset.DiameterCX,
		MaxIterations:         fractalPreset.MaxIterations,
		ColorPalette:          colorPreset.Palette,
		ColorPaletteRepeat:    fractalPreset.ColorPaletteRepeat,
		ColorPaletteLength:    fractalPreset.ColorPaletteLength,
		ColorPaletteReverse:   fractalPreset.ColorPaletteReverse,
		ColorPaletteHardStops: fractalPreset.ColorPaletteHardStops,
	}
	switch fractFunc {
	case FRACTAL_TYPE_MANDELBROT:
		return NewMandelbrotFractal(commonParams), nil
	case FRACTAL_TYPE_MANDELBROT3:
		return NewMandelbrot3Fractal(commonParams), nil
	case FRACTAL_TYPE_MANDELBROT4:
		return NewMandelbrot4Fractal(commonParams), nil
	case FRACTAL_TYPE_JULIA:
		return NewJuliaFractal(commonParams, fractalPreset.JuliaKr, fractalPreset.JuliaKi), nil
	default:
		return nil, errors.New("unknown fractal function")
	}
}

func NewFractalFromParams(fractFunct FractalType, commonFractParams CommonFractParams, juliaKr, juliaKi float64) (Fractal, error) {
	var fractal Fractal
	switch fractFunct {
	case "mandelbrot":
		fractal = NewMandelbrotFractal(commonFractParams)
		break
	case "mandelbrot3":
		fractal = NewMandelbrot3Fractal(commonFractParams)
		break
	case "mandelbrot4":
		fractal = NewMandelbrot4Fractal(commonFractParams)
		break
	case "julia":
		juliaKr := juliaKr
		juliaKi := juliaKi
		fractal = NewJuliaFractal(commonFractParams, juliaKr, juliaKi)
		break
	default:
		return nil, errors.New("unknown fractal function")
	}
	return fractal, nil
}
