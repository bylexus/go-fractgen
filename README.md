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
- [ ] WASM Frontend calc instead of Backend Server-side calc

## Build Notes

### CLI/Server

### WASM

Prep: Copy js runtime from go:

```sh
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" webroot/public/
```

Build wasm lib:

```sh
GOOS=js GOARCH=wasm go build -o webroot/public/main.wasm wasm/main.go
```

Initiate in index.html:

```html
    <script src="wasm_exec.js"></script>
    <script>
      const go = new Go();
      WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
          go.run(result.instance);
      });
	</script>
```