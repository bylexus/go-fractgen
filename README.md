# go-fractgen

A fractal generator written in Go, as cli, library and web app.

This is a work in progress. The main goal is to have the full functionality (and more) as my previous
project [JFractGen](https://github.com/bylexus/JFractGen).


### TODOs

#### General

- [ ] implement missing fractal functions (julia, mandelbrot ^ n, ...)

#### CLI

- [ ] create a cli app

#### Web app

- [ ] Caching of generated images
- [ ] Caching of fract values, to apply a new color scheme faster:
      If the user only changes the color scheme, the last fract values per user
	  should be saved. Needs some kind of session management.
- [x] zoom in/out with buttons
- [x] zoom in with double-click
- [ ] zoom in with drag a rectangle
- [ ] zoom in pinch zoom
- [ ] History / Undo stack
