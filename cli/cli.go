package cli

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"path"
	"strings"

	"github.com/bylexus/go-fract/lib"
	"github.com/bylexus/go-fract/web"
)

type ServeCmd struct {
	Listen      string `help:"Listen address / port to serve on." default:":8000"`
	PresetsFile string `help:"Path to presets file." type:"path"`
}

func (c *ServeCmd) Run(appContext *lib.AppContext) error {
	presets, err := lib.ReadPresetJson(c.PresetsFile, appContext.EmbeddedPresets)
	if err != nil {
		return err
	}
	server := web.NewWebServer(web.WebServerConfig{Addr: c.Listen}, presets.ColorPresets, presets.FractalPresets)

	fmt.Printf("Starting Webserver, listen on %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())

	return nil
}

type ImageCmd struct {
	Format         string          `help:"Format of the image to generate."`
	Width          int             `help:"Width of the image to generate, in pixels." default:"1920"`
	Height         int             `help:"Height of the image to generate, in pixels." default:"1200"`
	FractalPreset  string          `help:"Name of the fractal preset to use, e.g. '--fractal-preset=\"Mandelbrot Total\"'. Use in combination with --presets-file." default:"" `
	ColorPreset    string          `help:"Name of the color preset to use." default:"patchwork"`
	Function       lib.FractalType `help:"Fractal function to use." enum:"mandelbrot,julia,mandelbrot3,mandelbrot4" default:"mandelbrot"`
	PaletteRepeat  int             `help:"Number of times to repeat the palette." default:"1"`
	PaletteLength  int             `help:"Length of the palette." default:"-1"`
	PaletteReverse bool            `help:"Reverse the palette." default:"false"`
	CenterCX       string          `help:"Center CX(r)" default:"-0.7"`
	CenterCY       string          `help:"Center CY(i)" default:"0"`
	DiameterCX     string          `help:"Diameter CX(r)" default:"4"`
	JuliaKr        string          `help:"Julia Kr(r)" default:"-0.2"`
	JuliaKi        string          `help:"Julia Ki(i)" default:"0.8"`
	MaxIter        int             `help:"Maximum number of iterations." default:"100"`
	PresetsFile    string          `help:"Path to presets file." type:"path"`

	OutputPath string `arg:"" help:"Path to save the image to." type:"path" default:"image.jpg"`
}

func (c *ImageCmd) Run(appContext *lib.AppContext) error {
	presets, err := lib.ReadPresetJson(c.PresetsFile, appContext.EmbeddedPresets)
	if err != nil {
		return err
	}

	colorPreset, err := presets.ColorPresets.GetByIdent(c.ColorPreset)
	if err != nil {
		return err
	}

	var fractal lib.Fractal

	if c.Format == "" {
		ext := strings.ToLower(path.Ext(c.OutputPath))
		switch ext {
		case ".png":
			c.Format = "png"
		case ".jpg":
			fallthrough
		case ".jpeg":
			c.Format = "jpeg"
		default:
			return errors.New("unknown image format")
		}
	}

	if c.FractalPreset != "" {
		fractalPreset, err := presets.FractalPresets.GetByName(c.FractalPreset)
		if err != nil {
			return err
		}
		colorPreset, err = presets.ColorPresets.GetByIdent(fractalPreset.ColorPreset)
		if err != nil {
			return err
		}
		fractal, err = lib.NewFractalFromPresets(c.Width, c.Height, colorPreset, fractalPreset)
		if err != nil {
			return err
		}
		fmt.Printf("Using fractal preset: '%s', ignoring other fractal parameters.\n", c.FractalPreset)
	} else {
		centerCX, _, _ := big.ParseFloat(c.CenterCX, 10, lib.SYS_PRECISION, big.ToNearestEven)
		centerCY, _, _ := big.ParseFloat(c.CenterCY, 10, lib.SYS_PRECISION, big.ToNearestEven)
		diameterCX, _, _ := big.ParseFloat(c.DiameterCX, 10, lib.SYS_PRECISION, big.ToNearestEven)
		juliaKr, _, _ := big.ParseFloat(c.JuliaKr, 10, lib.SYS_PRECISION, big.ToNearestEven)
		juliaKi, _, _ := big.ParseFloat(c.JuliaKi, 10, lib.SYS_PRECISION, big.ToNearestEven)

		var commonFractParams = lib.CommonFractParams{
			ImageWidth:          c.Width,
			ImageHeight:         c.Height,
			CenterCX:            centerCX,
			CenterCY:            centerCY,
			DiameterCX:          diameterCX,
			MaxIterations:       c.MaxIter,
			ColorPalette:        colorPreset.Palette,
			ColorPaletteRepeat:  c.PaletteRepeat,
			ColorPaletteLength:  c.PaletteLength,
			ColorPaletteReverse: c.PaletteReverse,
		}
		fractal, err = lib.NewFractalFromParams(c.Function, commonFractParams, juliaKr, juliaKi)
		if err != nil {
			return err
		}
	}

	img := lib.CalcFractalImage(fractal)
	file, err := os.Create(c.OutputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	switch c.Format {
	case "png":
		err = img.EncodePng(file)
	case "jpeg":
		err = img.EncodeJpeg(file)
	default:
		return errors.New("unknown image format")
	}
	fmt.Printf("Image saved to %s\n", c.OutputPath)
	return err
}

type FlightCmd struct {
	Format          string          `help:"Format of the image to generate." enum:"png,jpeg,jpg" default:"jpeg"`
	Width           int             `help:"Width of the image to generate, in pixels." default:"720"`
	Height          int             `help:"Height of the image to generate, in pixels." default:"450"`
	ColorPreset     string          `help:"Name of the color preset to use." default:"patchwork"`
	Function        lib.FractalType `help:"Fractal function to use." enum:"mandelbrot,julia,mandelbrot3,mandelbrot4" default:"mandelbrot"`
	PaletteRepeat   int             `help:"Number of times to repeat the palette." default:"1"`
	PaletteLength   int             `help:"Length of the palette." default:"-1"`
	PaletteReverse  bool            `help:"Reverse the palette." default:"false"`
	StartCenterCX   string          `help:"Start Center CX(r)" default:"-0.7"`
	StartCenterCY   string          `help:"Start Center CY(i)" default:"0"`
	StartDiameterCX string          `help:"Start Diameter CX(r)" default:"4"`
	EndCenterCX     string          `help:"End Center CX(r)" default:"0.26954214666038734"`
	EndCenterCY     string          `help:"End Center CY(i)" default:"-0.00447479821741581"`
	EndDiameterCX   string          `help:"End Diameter CX(r)" default:"0.001220703125"`

	Duration int `help:"Duration of the flight, in Seconds." default:"10"`
	Fps      int `help:"Frames per second." default:"25"`

	JuliaKr     string `help:"Julia Kr(r)" default:"-0.2"`
	JuliaKi     string `help:"Julia Ki(i)" default:"0.8"`
	MaxIter     int    `help:"Maximum number of iterations." default:"800"`
	PresetsFile string `help:"Path to presets file." type:"path"`

	OutputFolder string `arg:"" help:"Folder to save the image to." type:"path" required:"true"`
}

// func (c *FlightCmd) Run(appContext *lib.AppContext) error {
// 	presets, err := lib.ReadPresetJson(c.PresetsFile, appContext.EmbeddedPresets)
// 	if err != nil {
// 		return err
// 	}

// 	colorPreset, err := presets.ColorPresets.GetByIdent(c.ColorPreset)
// 	if err != nil {
// 		return err
// 	}

// 	var fractal lib.Fractal
// 	startCenterCX, _, _ := big.ParseFloat(c.StartCenterCX, 10, lib.SYS_PRECISION, big.ToNearestEven)
// 	startCenterCY, _, _ := big.ParseFloat(c.StartCenterCY, 10, lib.SYS_PRECISION, big.ToNearestEven)
// 	startDiameterCX, _, _ := big.ParseFloat(c.StartDiameterCX, 10, lib.SYS_PRECISION, big.ToNearestEven)
// 	endCenterCX, _, _ := big.ParseFloat(c.EndCenterCX, 10, lib.SYS_PRECISION, big.ToNearestEven)
// 	endCenterCY, _, _ := big.ParseFloat(c.EndCenterCY, 10, lib.SYS_PRECISION, big.ToNearestEven)
// 	endDiameterCX, _, _ := big.ParseFloat(c.EndDiameterCX, 10, lib.SYS_PRECISION, big.ToNearestEven)
// 	juliaKr, _, _ := big.ParseFloat(c.JuliaKr, 10, lib.SYS_PRECISION, big.ToNearestEven)
// 	juliaKi, _, _ := big.ParseFloat(c.JuliaKi, 10, lib.SYS_PRECISION, big.ToNearestEven)

// 	var commonFractParams = lib.CommonFractParams{
// 		ImageWidth:          c.Width,
// 		ImageHeight:         c.Height,
// 		CenterCX:            startCenterCX,
// 		CenterCY:            startCenterCY,
// 		DiameterCX:          startDiameterCX,
// 		MaxIterations:       c.MaxIter,
// 		ColorPalette:        colorPreset.Palette,
// 		ColorPaletteRepeat:  c.PaletteRepeat,
// 		ColorPaletteLength:  c.PaletteLength,
// 		ColorPaletteReverse: c.PaletteReverse,
// 	}

// 	nrOfImages := c.Duration * c.Fps

// 	err = os.MkdirAll(c.OutputFolder, os.ModePerm)
// 	if err != nil {
// 		return err
// 	}
// 	deltaX := new(big.Float).Sub(endCenterCX, startCenterCX)

// 	deltaY := new(big.Float).Sub(endCenterCY, startCenterCY)

// 	deltaDiameter := new(big.Float).Sub(endDiameterCX, startDiameterCX)

// 	/*
// 		We cannot simply increase the diameter by the same amount every step: as we dive deeper into
// 		the fractal, we are "nearer" to the end point: This means that the same amount of diameter increase
// 		looks like we will get faster and faster.

// 		To counteract this, we need to increase the diameter by an amount that scales with the current
// 		distance from the end point. Or, we increase/decrease the diameter by a percentage amount of the actual value.

// 		This seems to be a similar problem to the "Zineszins" problem:

// 		Kn = K0 * (1 + p)^n

// 		where:
// 		- K0 is the initial value
// 		- p is the percentage increase/decrease, from 0 to 1
// 		- n is the number of steps
// 		- Kn is the final value

// 		We have:
// 		- the initial value K0: c.StartDiameterCX
// 		- the final value Kn: c.EndDiameterCX
// 		- the number of steps (n, nrOfImages)
// 		- We want: The value of p:
// 			p = (nth root of (Kn / K0)) + 1

// 			Example:

// 			nrOfImages = 80
// 			K0 = 4
// 			Kn = 0.01

// 			p = (Kn/K(0))^(1/80) - 1
// 			p = (0.01/4)^(1/80) - 1
// 			p = 1.0000000000000002
// 	*/

// 	// diameter percentage to increase per image:
// 	// p = (nth root of (endDiameter / startDiameter)) + 1
// 	var percIncrease float64 = math.Pow(float64(c.EndDiameterCX)/float64(c.StartDiameterCX), 1.0/float64(nrOfImages)) - 1.0

// 	// i := 0
// 	inc := big.NewFloat(1.0)
// 	inc.Add(inc, big.NewFloat(percIncrease))

// 	actDiameterCX := new(big.Float).Copy(startDiameterCX)
// 	// fmt.Printf("Start Diameter CX: %v, End Diameter CX: %v, Inc Diameter CX: %v\n", startDiameterCX, endDiameterCX, inc)

// 	// for actDiameterCX.Cmp(endDiameterCX) == dir
// 	for i := 0; i <= nrOfImages; i++ {
// 		total := new(big.Float).Copy(deltaDiameter)
// 		diff := new(big.Float).Sub(actDiameterCX, startDiameterCX)
// 		percDone := diff.Quo(diff, total).Abs(diff)

// 		// newCX = deltaX * percentage + startCenterCX
// 		newCX := new(big.Float).Copy(deltaX)
// 		newCX.Mul(newCX, percDone)
// 		newCX.Add(newCX, startCenterCX)

// 		// newCY = deltaY * percentage + startCenterCY
// 		newCY := new(big.Float).Copy(deltaY)
// 		newCY.Mul(newCY, percDone)
// 		newCY.Add(newCY, startCenterCY)

// 		commonFractParams.CenterCX, _ = newCX.Float64()
// 		commonFractParams.CenterCY, _ = newCY.Float64()
// 		commonFractParams.DiameterCX, _ = actDiameterCX.Mul(actDiameterCX, inc).Float64()

// 		fractal, err = lib.NewFractalFromParams(c.Function, commonFractParams, c.JuliaKr, c.JuliaKi)
// 		if err != nil {
// 			return err
// 		}
// 		filename := fmt.Sprintf("%s/%08d.%s", c.OutputFolder, i, c.Format)

// 		img := lib.CalcFractalImage(fractal)
// 		file, err := os.Create(filename)
// 		if err != nil {
// 			return err
// 		}
// 		defer file.Close()

// 		switch c.Format {
// 		case "png":
// 			err = img.EncodePng(file)
// 		case "jpeg":
// 			fallthrough
// 		case "jpg":
// 			err = img.EncodeJpeg(file)
// 		default:
// 			return errors.New("unknown image format")
// 		}
// 		p := float64(i) / float64(nrOfImages)
// 		fmt.Printf("%0.1f%% Image saved to %s\n", p*100, filename)

// 	}
// 	fmt.Printf("End diameter cx: %v\n", commonFractParams.DiameterCX)

// 	// now, convert images with ffmpeg:
// 	//ffmpeg -framerate 25 -pattern_type glob -i '*.jpeg' -c:v libx264 -pix_fmt yuv420p out.mp4

// 	return err
// }

type Cli struct {
	Serve  ServeCmd  `cmd:"" help:"Start the web server."`
	Image  ImageCmd  `cmd:"" help:"Generate a single image."`
	Flight FlightCmd `cmd:"" help:"Generate a flight through a fractal: generate a series of images from a start point to an end point."`
}
