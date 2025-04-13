# go-fractgen

A fractal generator written in Go, as cli, library and web app.

This is a work in progress. The main goal is to have the full functionality (and more) as my previous
project [JFractGen](https://github.com/bylexus/JFractGen).

## Features

(TODO)

- Using OpenLayers in frontend to navigate on a virtual Fractal WMTS tile grid
- renders image tiles in the golang backend


### TODOs

#### General

- [ ] implement missing fractal functions (julia, mandelbrot ^ n, ...)

#### CLI

- [ ] create a cli app

#### Web app

- [x] Treat the fractal plane as grid: Create a grid-based approach to generate and cache the images,
using a Mapping library like OpenLayers to render the tiles. This allows for easy pre-generation and caching
of the fractal images as tiles.
- [ ] Caching of generated images, using wmts tiles from above, and the necessary metadata as key
- [x] zoom in/out with buttons
- [x] zoom in with double-click
- [x] zoom in with drag a rectangle
- [x] zoom in pinch zoom
- [ ] History / Undo stack
