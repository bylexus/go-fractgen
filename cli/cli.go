package cli

import (
	"errors"
	"fmt"
	"log"
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

type Cli struct {
	Serve ServeCmd `cmd:"" help:"Start the web server."`
	Image ImageCmd `cmd:"" help:"Generate a single image."`
}
