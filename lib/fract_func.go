package lib

import (
	"errors"
	"math/big"
	"runtime"
	"strings"

	"github.com/bylexus/go-stdlib/ethreads"
)

// max amount = 256 gives smoother colors around the fractal:
// const MAX_ABS_SQUARE_AMOUNT float64 = 4
// const MAX_ABS_SQUARE_AMOUNT float64 = 256

type FractalType string

const (
	FRACTAL_TYPE_MANDELBROT  = "mandelbrot"
	FRACTAL_TYPE_MANDELBROT3 = "mandelbrot3"
	FRACTAL_TYPE_MANDELBROT4 = "mandelbrot4"
	FRACTAL_TYPE_JULIA       = "julia"
)

func FractalTypeFromString(s string) (FractalType, error) {
	switch strings.ToLower(s) {
	case "mandelbrot":
		return FRACTAL_TYPE_MANDELBROT, nil
	case "mandelbrot3":
		return FRACTAL_TYPE_MANDELBROT3, nil
	case "mandelbrot4":
		return FRACTAL_TYPE_MANDELBROT4, nil
	case "julia":
		return FRACTAL_TYPE_JULIA, nil
	default:
		return "", errors.New("unknown fractal function")
	}
}

type Fractal interface {
	CreatePixelCalcFunc(pixX, pixY int, img *FractImage) ethreads.JobFn
	ImageWidth() int
	ImageHeight() int
}

type FractFunctionResult struct {
	Iterations   int
	BailoutValue *big.Float
}

type CommonFractParams struct {
	// MaxAbsSquareAmount *big.Float
	MaxIterations int

	CenterCX   *big.Float
	CenterCY   *big.Float
	DiameterCX *big.Float

	ImageWidth  int
	ImageHeight int

	SmoothColors        bool
	ColorPaletteLength  int
	ColorPalette        ColorPalette
	ColorPaletteRepeat  int
	ColorPaletteReverse bool

	// calculaed during initialization:
	aspect *big.Float
	minCX  *big.Float
	maxCX  *big.Float
	minCY  *big.Float
	maxCY  *big.Float
}

func (f CommonFractParams) PixelToFractal(x, y int) (cx, cy *big.Float) {
	// y axis is inverted in image and fractal space:
	y = f.ImageHeight - y
	percentageX := big.NewFloat(float64(x) / float64(f.ImageWidth)).SetPrec(SYS_PRECISION)
	percentageY := big.NewFloat(float64(y) / float64(f.ImageHeight)).SetPrec(SYS_PRECISION)
	diameterCX := new(big.Float).Copy(f.DiameterCX)

	diameterCY := new(big.Float).Copy(f.maxCY)
	diameterCY.Sub(diameterCY, f.minCY)

	// cx = f.minCX + (f.DiameterCX)*(float64(x)/float64(f.ImageWidth))
	cx = new(big.Float).Copy(f.minCX)
	cx.Add(cx, percentageX.Mul(diameterCX, percentageX))

	// cy = f.minCY + (f.maxCY-f.minCY)*(float64(y)/float64(f.ImageHeight))
	cy = new(big.Float).Copy(f.minCY)
	cy.Add(cy, percentageY.Mul(diameterCY, percentageY))
	return cx, cy
}

func initializeFractParams(commonFractParams CommonFractParams) CommonFractParams {
	// var aspect, fract_width, fract_heigth float64
	var aspect, fract_width, fract_heigth, centerCX, centerCY *big.Float
	var bigWidth = big.NewFloat(float64(commonFractParams.ImageWidth)).SetPrec(SYS_PRECISION)
	var bigHeight = big.NewFloat(float64(commonFractParams.ImageHeight)).SetPrec(SYS_PRECISION)

	aspect = new(big.Float).SetPrec(SYS_PRECISION).Quo(bigWidth, bigHeight)
	// aspect = float64(commonFractParams.ImageWidth) / float64(commonFractParams.ImageHeight)
	fract_width = new(big.Float).Copy(commonFractParams.DiameterCX)
	fract_heigth = new(big.Float).SetPrec(SYS_PRECISION).Quo(fract_width, aspect)
	centerCX = new(big.Float).Copy(commonFractParams.CenterCX)
	centerCY = new(big.Float).Copy(commonFractParams.CenterCY)

	// var min_cx float64 = commonFractParams.CenterCX - (fract_width / 2.0)
	var min_cx = new(big.Float).SetPrec(SYS_PRECISION).Sub(centerCX, (new(big.Float).SetPrec(SYS_PRECISION).Quo(fract_width, BIG_2)))
	// var max_cx float64 = min_cx + fract_width
	var max_cx = new(big.Float).SetPrec(SYS_PRECISION).Add(min_cx, fract_width)
	// var min_cy float64 = commonFractParams.CenterCY - (fract_heigth / 2.0)
	var min_cy = new(big.Float).SetPrec(SYS_PRECISION).Sub(centerCY, (new(big.Float).SetPrec(SYS_PRECISION).Quo(fract_heigth, BIG_2)))
	// var max_cy float64 = min_cy + fract_heigth
	var max_cy = new(big.Float).SetPrec(SYS_PRECISION).Add(min_cy, fract_heigth)

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
	// commonFractParams.MaxAbsSquareAmount = big.NewFloat(MAX_ABS_SQUARE_AMOUNT).SetPrec(SYS_PRECISION)

	return commonFractParams
}

func CalcFractalImage(f Fractal) *FractImage {
	tp := ethreads.NewThreadPool(runtime.NumCPU()*2, nil)
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
		ImageWidth:          width,
		ImageHeight:         height,
		CenterCX:            fractalPreset.CenterCX,
		CenterCY:            fractalPreset.CenterCY,
		DiameterCX:          fractalPreset.DiameterCX,
		MaxIterations:       fractalPreset.MaxIterations,
		ColorPalette:        colorPreset.Palette,
		ColorPaletteRepeat:  fractalPreset.ColorPaletteRepeat,
		ColorPaletteLength:  fractalPreset.ColorPaletteLength,
		ColorPaletteReverse: fractalPreset.ColorPaletteReverse,
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

func NewFractalFromParams(fractFunct FractalType, commonFractParams CommonFractParams, juliaKr, juliaKi *big.Float) (Fractal, error) {
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
