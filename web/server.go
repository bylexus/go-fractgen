package web

import (
	"encoding/json"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"github.com/bylexus/go-fract/lib"
	"github.com/bylexus/go-stdlib/ethreads"
)

type WebServer struct {
	http.Server

	threadPool *ethreads.ThreadPool

	colorPresets   lib.ColorPresets
	fractalPresets lib.FractalPresets
}

func NewWebServer(colorPresets lib.ColorPresets, fractalPresets lib.FractalPresets) *WebServer {
	threadPool := ethreads.NewThreadPool(runtime.NumCPU(), nil)
	threadPool.Start()
	server := &WebServer{
		threadPool:     &threadPool,
		colorPresets:   colorPresets,
		fractalPresets: fractalPresets,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/fractal-image.jpg", server.handleFractalImage)
	mux.HandleFunc("/wmts", server.handleWmtsRequest)
	mux.HandleFunc("/presets.json", server.handlePresetsJson)
	mux.Handle("/", http.FileServer(http.Dir("webroot")))

	listenAddr := ":8000"
	var handler http.Handler = NewLogHandler(NewCorsHandler(mux))
	server.Server = http.Server{
		Addr: listenAddr, Handler: handler,
	}

	return server
}

func (s *WebServer) handleFractalImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/jpeg")
	width, _ := strconv.Atoi(r.URL.Query().Get("width"))
	height, _ := strconv.Atoi(r.URL.Query().Get("height"))
	maxIterations, _ := strconv.Atoi(r.URL.Query().Get("maxIterations"))
	iterFunc := strings.ToLower(r.URL.Query().Get("iterFunc"))
	centerCX, _ := strconv.ParseFloat(r.URL.Query().Get("centerCX"), 64)
	centerCY, _ := strconv.ParseFloat(r.URL.Query().Get("centerCY"), 64)
	diameterCX, _ := strconv.ParseFloat(r.URL.Query().Get("diameterCX"), 64)
	colorPresetParam := r.URL.Query().Get("colorPreset")
	colorPaletteRepeat, _ := strconv.Atoi(r.URL.Query().Get("colorPaletteRepeat"))

	colorPreset, err := s.colorPresets.GetByName(colorPresetParam)
	if err != nil {
		colorPreset = s.colorPresets[0]
	}

	var fractal lib.Fractal
	switch iterFunc {
	case "mandelbrot":
		fractal = lib.NewMandelbrotFractal(
			width, height,
			centerCX, centerCY, diameterCX,
			maxIterations, colorPreset.Palette,
			colorPaletteRepeat,
		)
		break
	case "mandelbrot3":
		fractal = lib.NewMandelbrot3Fractal(
			width, height,
			centerCX, centerCY, diameterCX,
			maxIterations, colorPreset.Palette,
			colorPaletteRepeat,
		)
		break
	case "mandelbrot4":
		fractal = lib.NewMandelbrot4Fractal(
			width, height,
			centerCX, centerCY, diameterCX,
			maxIterations, colorPreset.Palette,
			colorPaletteRepeat,
		)
		break
	case "julia":
		juliaKr, _ := strconv.ParseFloat(r.URL.Query().Get("juliaKr"), 64)
		juliaKi, _ := strconv.ParseFloat(r.URL.Query().Get("juliaKi"), 64)
		fractal = lib.NewJuliaFractal(
			width, height,
			centerCX, centerCY, diameterCX,
			maxIterations, colorPreset.Palette,
			colorPaletteRepeat,
			juliaKr, juliaKi,
		)
		break
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	img := lib.CalcFractalImage(s.threadPool, fractal)
	w.WriteHeader(http.StatusOK)
	img.EncodeJpeg(w)
}

func (s *WebServer) handleWmtsRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/jpeg")
	zoomLevel, _ := strconv.Atoi(r.URL.Query().Get("TileMatrix"))
	if zoomLevel < 0 || zoomLevel > 50 {
		zoomLevel = 0
	}

	iterFunc := strings.ToLower(r.URL.Query().Get("iterFunc"))

	tileX, _ := strconv.Atoi(r.URL.Query().Get("TileCol"))
	tileY, _ := strconv.Atoi(r.URL.Query().Get("TileRow"))

	originX := -1.7
	originY := -1.0

	tileWidthPixels, _ := strconv.Atoi(r.URL.Query().Get("tileWidthPixels"))
	if tileWidthPixels == 0 {
		tileWidthPixels = 256
	}
	tileWidthFractal, _ := strconv.ParseFloat(r.URL.Query().Get("tileWidthFractal"), 64)
	if tileWidthFractal == 0 {
		tileWidthFractal = 1
	}

	centerCX := originX + float64(tileX)*tileWidthFractal + (tileWidthFractal / 2)
	centerCY := originY + float64(-1*tileY)*tileWidthFractal - (tileWidthFractal / 2)

	maxIterations, _ := strconv.Atoi(r.URL.Query().Get("maxIterations"))
	if maxIterations == 0 {
		maxIterations = 50
	}
	colorPresetParam := r.URL.Query().Get("colorPreset")
	if colorPresetParam == "" {
		colorPresetParam = "Patchwork"
	}
	colorPaletteRepeat, _ := strconv.Atoi(r.URL.Query().Get("colorPaletteRepeat"))
	if colorPaletteRepeat == 0 {
		colorPaletteRepeat = 1
	}

	colorPreset, err := s.colorPresets.GetByName(colorPresetParam)
	if err != nil {
		colorPreset = s.colorPresets[0]
	}

	var fractal lib.Fractal
	switch iterFunc {
	case "mandelbrot":
		fractal = lib.NewMandelbrotFractal(
			tileWidthPixels, tileWidthPixels,
			centerCX, centerCY, tileWidthFractal,
			maxIterations, colorPreset.Palette,
			colorPaletteRepeat,
		)
		break
	case "mandelbrot3":
		fractal = lib.NewMandelbrot3Fractal(
			tileWidthPixels, tileWidthPixels,
			centerCX, centerCY, tileWidthFractal,
			maxIterations, colorPreset.Palette,
			colorPaletteRepeat,
		)
		break
	case "mandelbrot4":
		fractal = lib.NewMandelbrot4Fractal(
			tileWidthPixels, tileWidthPixels,
			centerCX, centerCY, tileWidthFractal,
			maxIterations, colorPreset.Palette,
			colorPaletteRepeat,
		)
		break
	case "julia":
		juliaKr, _ := strconv.ParseFloat(r.URL.Query().Get("juliaKr"), 64)
		juliaKi, _ := strconv.ParseFloat(r.URL.Query().Get("juliaKi"), 64)
		fractal = lib.NewJuliaFractal(
			tileWidthPixels, tileWidthPixels,
			centerCX, centerCY, tileWidthFractal,
			maxIterations, colorPreset.Palette,
			colorPaletteRepeat,
			juliaKr, juliaKi,
		)
		break
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	img := lib.CalcFractalImage(s.threadPool, fractal)
	w.Header().Set("Cache-Control", "public, max-age=15552000;")
	w.WriteHeader(http.StatusOK)
	img.EncodeJpeg(w)
}

func (s *WebServer) handlePresetsJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	presets := lib.Presets{
		ColorPresets:   s.colorPresets,
		FractalPresets: s.fractalPresets,
	}

	jsonStream, err := json.Marshal(presets)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStream)
}
