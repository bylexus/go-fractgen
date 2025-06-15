/**
 * Fractal Generator written by Alexander Schenkel, www.alexi.ch
 *
 * This program implements a fractal generator to generate Mandelbrot and
 * Julia set fractals. It offers the following UIs:
 *
 * - a Web server (run with the `serve` command)
 * - a command line interface to generate single images (run with the `image` command)
 * - a command line interface to generate an image series / flight through a fractal (run with the `flight` command)
 *
 * Build with:
 *
 *      go build -o fractgen
 *
 * (c) 2012-2025 Alexander Schenkel
 */
package main

import (
	"embed"
	_ "embed"
	"io/fs"

	"github.com/alecthomas/kong"
	"github.com/bylexus/go-fract/cli"
	"github.com/bylexus/go-fract/lib"
)

//go:embed presets.json
var presets []byte

//go:embed webroot/dist
var webrootDist embed.FS

func main() {

	webrootFS, err := fs.Sub(webrootDist, "webroot/dist")

	if err != nil {
		panic(err)
	}
	cli := cli.Cli{}
	ctx := kong.Parse(&cli)
	err = ctx.Run(&lib.AppContext{EmbeddedPresets: presets, WebrootFS: webrootFS})
	ctx.FatalIfErrorf(err)
}
