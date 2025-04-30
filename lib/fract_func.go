package lib

import (
	"errors"
	"runtime"

	"github.com/bylexus/go-stdlib/ethreads"
)

// max amount = 256 gives smoother colors around the fractal:
// const MAX_ABS_SQUARE_AMOUNT float64 = 4
const MAX_ABS_SQUARE_AMOUNT float64 = 256

type FractalType int

const (
	FRACTAL_TYPE_MANDELBROT = iota
	FRACTAL_TYPE_MANDELBROT3
	FRACTAL_TYPE_MANDELBROT4
	FRACTAL_TYPE_JULIA
)

type Fractal interface {
	CreatePixelCalcFunc(pixX, pixY int, img *FractImage) ethreads.JobFn
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

	SmoothColors       bool
	ColorPaletteLength int
	ColorPalette       ColorPalette
	ColorPaletteRepeat int

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
	var aspect, fract_width, fract_heigth float64

	aspect = float64(commonFractParams.ImageWidth) / float64(commonFractParams.ImageHeight)
	fract_width = commonFractParams.DiameterCX
	fract_heigth = commonFractParams.DiameterCX / aspect

	var min_cx float64 = commonFractParams.CenterCX - (fract_width / 2.0)
	var max_cx float64 = min_cx + fract_width
	var min_cy float64 = commonFractParams.CenterCY - (fract_heigth / 2.0)
	var max_cy float64 = min_cy + fract_heigth

	if commonFractParams.ColorPaletteRepeat <= 0 {
		commonFractParams.ColorPaletteRepeat = 1
	}
	if commonFractParams.ColorPaletteLength <= 0 {
		commonFractParams.ColorPaletteLength = -1
	}

	commonFractParams.SmoothColors = true

	// Calculated during initialization:
	commonFractParams.aspect = aspect
	commonFractParams.minCX = min_cx
	commonFractParams.maxCX = max_cx
	commonFractParams.minCY = min_cy
	commonFractParams.maxCY = max_cy
	commonFractParams.MaxAbsSquareAmount = MAX_ABS_SQUARE_AMOUNT

	return commonFractParams

}

func CalcFractalImage(f Fractal) *FractImage {
	tp := ethreads.NewThreadPool(runtime.NumCPU(), nil)
	tp.Start()

	img := NewFractImage(f.ImageWidth(), f.ImageHeight())

	for y := 0; y < f.ImageHeight(); y++ {
		for x := 0; x < f.ImageWidth(); x++ {
			tp.AddJobFn(f.CreatePixelCalcFunc(x, y, img))
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
		ImageWidth:         width,
		ImageHeight:        height,
		CenterCX:           fractalPreset.CenterCX,
		CenterCY:           fractalPreset.CenterCY,
		DiameterCX:         fractalPreset.DiameterCX,
		MaxIterations:      fractalPreset.MaxIterations,
		ColorPalette:       colorPreset.Palette,
		ColorPaletteRepeat: fractalPreset.ColorPaletteRepeat,
		ColorPaletteLength: fractalPreset.ColorPaletteLength,
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
