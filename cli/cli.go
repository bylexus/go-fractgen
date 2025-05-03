package cli

import (
	"errors"
	"fmt"
	"log"
	"math"
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
	CenterCX       float64         `help:"Center CX(r)" default:"-0.7"`
	CenterCY       float64         `help:"Center CY(i)" default:"0"`
	DiameterCX     float64         `help:"Diameter CX(r)" default:"4"`
	JuliaKr        float64         `help:"Julia Kr(r)" default:"-0.2"`
	JuliaKi        float64         `help:"Julia Ki(i)" default:"0.8"`
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
		fmt.Println("Note: Using a fractal preset, ignoring other fractal parameters.")
	} else {
		var commonFractParams = lib.CommonFractParams{
			ImageWidth:          c.Width,
			ImageHeight:         c.Height,
			CenterCX:            c.CenterCX,
			CenterCY:            c.CenterCY,
			DiameterCX:          c.DiameterCX,
			MaxIterations:       c.MaxIter,
			ColorPalette:        colorPreset.Palette,
			ColorPaletteRepeat:  c.PaletteRepeat,
			ColorPaletteLength:  c.PaletteLength,
			ColorPaletteReverse: c.PaletteReverse,
		}
		fractal, err = lib.NewFractalFromParams(c.Function, commonFractParams, c.JuliaKr, c.JuliaKi)
		if err != nil {
			return err
		}
	}

	if c.Format == "" {
		ext := strings.ToLower(path.Ext(c.OutputPath))
		switch ext {
		case ".png":
			c.Format = "png"
		case ".jpg":
		case ".jpeg":
			c.Format = "jpeg"
		default:
			return errors.New("unknown image format")
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
	StartCenterCX   float64         `help:"Start Center CX(r)" default:"-0.7"`
	StartCenterCY   float64         `help:"Start Center CY(i)" default:"0"`
	StartDiameterCX float64         `help:"Start Diameter CX(r)" default:"4"`
	EndCenterCX     float64         `help:"End Center CX(r)" default:"0.26954214666038734"`
	EndCenterCY     float64         `help:"End Center CY(i)" default:"-0.00447479821741581"`
	EndDiameterCX   float64         `help:"End Diameter CX(r)" default:"0.001220703125"`

	Duration int `help:"Duration of the flight, in Seconds." default:"10"`
	Fps      int `help:"Frames per second." default:"25"`

	JuliaKr     float64 `help:"Julia Kr(r)" default:"-0.2"`
	JuliaKi     float64 `help:"Julia Ki(i)" default:"0.8"`
	MaxIter     int     `help:"Maximum number of iterations." default:"800"`
	PresetsFile string  `help:"Path to presets file." type:"path"`

	OutputFolder string `arg:"" help:"Folder to save the image to." type:"path" required:"true"`
}

func (c *FlightCmd) Run(appContext *lib.AppContext) error {
	presets, err := lib.ReadPresetJson(c.PresetsFile, appContext.EmbeddedPresets)
	if err != nil {
		return err
	}

	colorPreset, err := presets.ColorPresets.GetByIdent(c.ColorPreset)
	if err != nil {
		return err
	}

	var fractal lib.Fractal

	var commonFractParams = lib.CommonFractParams{
		ImageWidth:          c.Width,
		ImageHeight:         c.Height,
		CenterCX:            c.StartCenterCX,
		CenterCY:            c.StartCenterCY,
		DiameterCX:          c.StartDiameterCX,
		MaxIterations:       c.MaxIter,
		ColorPalette:        colorPreset.Palette,
		ColorPaletteRepeat:  c.PaletteRepeat,
		ColorPaletteLength:  c.PaletteLength,
		ColorPaletteReverse: c.PaletteReverse,
	}

	nrOfImages := c.Duration * c.Fps

	err = os.MkdirAll(c.OutputFolder, os.ModePerm)
	if err != nil {
		return err
	}
	startCenterCX := big.NewFloat(c.StartCenterCX)
	endCenterCX := big.NewFloat(c.EndCenterCX)
	deltaX := big.NewFloat(0.0).Sub(endCenterCX, startCenterCX)
	// incX := big.NewFloat(0.0).Quo(deltaX, big.NewFloat(float64(nrOfImages)))

	startCenterCY := big.NewFloat(c.StartCenterCY)
	endCenterCY := big.NewFloat(c.EndCenterCY)
	deltaY := big.NewFloat(0.0).Sub(endCenterCY, startCenterCY)
	// incY := big.NewFloat(0.0).Quo(deltaY, big.NewFloat(float64(nrOfImages)))

	startDiameterCX := big.NewFloat(c.StartDiameterCX)
	endDiameterCX := big.NewFloat(c.EndDiameterCX)
	deltaDiameter := big.NewFloat(0.0).Sub(endDiameterCX, startDiameterCX)
	// incDiameter := big.NewFloat(0.0).Quo(deltaDiameter, big.NewFloat(float64(nrOfImages)))
	// incDiameter.Mul(incDiameter, big.NewFloat(2.0))

	// fmt.Printf("Start Center CX: %v, End Center CX: %v, Inc Center CX: %v\n", startCenterCX, endCenterCX, incX)

	/*
		We cannot simply increase the diameter by the same amount every step: as we dive deeper into
		the fractal, we are "nearer" to the end point: This means that the same amount of diameter increase
		looks like we will get faster and faster.

		To counteract this, we need to increase the diameter by an amount that scales with the current
		distance from the end point. Or, we increase/decrease the diameter by a percentage amount of the actual value.

		This seems to be a similar problem to the "Zineszins" problem:

		Kn = K0 * (1 + p)^n

		where:
		- K0 is the initial value
		- p is the percentage increase/decrease, from 0 to 1
		- n is the number of steps
		- Kn is the final value

		We have:
		- the initial value K0: c.StartDiameterCX
		- the final value Kn: c.EndDiameterCX
		- the number of steps (n, nrOfImages)
		- We want: The value of p:
			p = (nth root of (Kn / K0)) + 1

			Example:

			nrOfImages = 80
			K0 = 4
			Kn = 0.01

			p = (Kn/K(0))^(1/80) - 1
			p = (0.01/4)^(1/80) - 1
			p = 1.0000000000000002



	*/

	// dir := startDiameterCX.Cmp(endDiameterCX)

	// diameter percentage to increase per image:
	// p = (nth root of (endDiameter / startDiameter)) + 1
	var percIncrease float64 = math.Pow(float64(c.EndDiameterCX)/float64(c.StartDiameterCX), 1.0/float64(nrOfImages)) - 1.0

	// i := 0
	inc := big.NewFloat(1.0)
	inc.Add(inc, big.NewFloat(percIncrease))

	actDiameterCX := big.NewFloat(0.0).Copy(startDiameterCX)
	// fmt.Printf("Start Diameter CX: %v, End Diameter CX: %v, Inc Diameter CX: %v\n", startDiameterCX, endDiameterCX, inc)

	// for actDiameterCX.Cmp(endDiameterCX) == dir {
	for i := 0; i <= nrOfImages; i++ {
		total := big.NewFloat(0.0).Copy(deltaDiameter)
		diff := big.NewFloat(0.0).Sub(actDiameterCX, startDiameterCX)
		percDone := diff.Quo(diff, total).Abs(diff)

		// commonFractParams.CenterCX = c.StartCenterCX + (c.EndCenterCX-c.StartCenterCX)*float64(i+1)/float64(nrOfImages)
		// commonFractParams.CenterCY = c.StartCenterCY + (c.EndCenterCY-c.StartCenterCY)*float64(i+1)/float64(nrOfImages)
		// commonFractParams.DiameterCX = c.StartDiameterCX + (c.EndDiameterCX-c.StartDiameterCX)*float64(i+1)/float64(nrOfImages)

		// prev := big.NewFloat(0.0).Copy(actDiameterCX)
		if i == 0 {
			// newcx = startCenterCX + (deltaX)*percentage

			commonFractParams.CenterCX, _ = startCenterCX.Float64()
			commonFractParams.CenterCY, _ = startCenterCY.Float64()
			commonFractParams.DiameterCX, _ = actDiameterCX.Float64()

		} else {
			newCX := big.NewFloat(0.0).Copy(deltaX)
			newCX.Mul(newCX, percDone)
			newCX.Add(newCX, startCenterCX)

			newCY := big.NewFloat(0.0).Copy(deltaY)
			newCY.Mul(newCY, percDone)
			newCY.Add(newCY, startCenterCY)
			// commonFractParams.CenterCX, _ = startCenterCX.Add(startCenterCX, incX).Float64()
			// commonFractParams.CenterCY, _ = startCenterCY.Add(startCenterCY, incY).Float64()
			// commonFractParams.DiameterCX, _ = actDiameterCX.Add(actDiameterCX, incDiameter).Float64()
			commonFractParams.CenterCX, _ = newCX.Float64()
			commonFractParams.CenterCY, _ = newCY.Float64()
			commonFractParams.DiameterCX, _ = actDiameterCX.Mul(actDiameterCX, inc).Float64()

		}
		// diff := big.NewFloat(0.0).Sub(actDiameterCX, prev)
		// fmt.Printf("Start diameter cx: %v, End diameter cx: %v, Diff: %v\n", prev, actDiameterCX, diff)

		fractal, err = lib.NewFractalFromParams(c.Function, commonFractParams, c.JuliaKr, c.JuliaKi)
		if err != nil {
			return err
		}
		filename := fmt.Sprintf("%s/%08d.%s", c.OutputFolder, i, c.Format)

		img := lib.CalcFractalImage(fractal)
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		switch c.Format {
		case "png":
			err = img.EncodePng(file)
		case "jpeg":
			fallthrough
		case "jpg":
			err = img.EncodeJpeg(file)
		default:
			return errors.New("unknown image format")
		}
		p := float64(i) / float64(nrOfImages)
		fmt.Printf("%0.1f%% Image saved to %s\n", p*100, filename)

	}
	fmt.Printf("End diameter cx: %v\n", commonFractParams.DiameterCX)

	// now, convert images with ffmpeg:
	//ffmpeg -framerate 25 -pattern_type glob -i '*.jpeg' -c:v libx264 -pix_fmt yuv420p out.mp4

	return err
}

type Cli struct {
	Serve  ServeCmd  `cmd:"" help:"Start the web server."`
	Image  ImageCmd  `cmd:"" help:"Generate a single image."`
	Flight FlightCmd `cmd:"" help:"Generate a flight through a fractal: generate a series of images from a start point to an end point."`
}
