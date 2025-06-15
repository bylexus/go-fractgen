package web

import (
	"encoding/json"
	"fmt"
	"image"
	"net/http"
	"strconv"
	"strings"

	"github.com/bylexus/go-fract/lib"
)

type WebServerConfig struct {
	Addr string
}

type WebServer struct {
	http.Server

	colorPresets   lib.ColorPresets
	fractalPresets lib.FractalPresets
}

func NewWebServer(conf WebServerConfig, colorPresets lib.ColorPresets, fractalPresets lib.FractalPresets) *WebServer {
	server := &WebServer{
		colorPresets:   colorPresets,
		fractalPresets: fractalPresets,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/fractal-image/{format}", server.handleFractalImage)
	mux.HandleFunc("/paletteViewer", server.handlePaletteViewer)
	mux.HandleFunc("/wmts", server.handleWmtsRequest)
	mux.HandleFunc("/presets.json", server.handlePresetsJson)
	mux.HandleFunc("/ping", server.handlePing)
	mux.Handle("/", http.FileServer(http.Dir("webroot")))

	listenAddr := conf.Addr
	var handler http.Handler = NewLogHandler(NewCorsHandler(mux))
	server.Server = http.Server{
		Addr: listenAddr, Handler: handler,
	}

	return server
}

func (s *WebServer) handleFractalImage(w http.ResponseWriter, r *http.Request) {
	format := strings.ToLower(r.PathValue("format"))
	width, _ := strconv.Atoi(r.URL.Query().Get("width"))
	height, _ := strconv.Atoi(r.URL.Query().Get("height"))
	maxIterations, _ := strconv.Atoi(r.URL.Query().Get("maxIterations"))
	iterFunc := strings.ToLower(r.URL.Query().Get("iterFunc"))
	centerCX, _ := strconv.ParseFloat(r.URL.Query().Get("centerCX"), 64)
	centerCY, _ := strconv.ParseFloat(r.URL.Query().Get("centerCY"), 64)
	diameterCX, _ := strconv.ParseFloat(r.URL.Query().Get("diameterCX"), 64)
	colorPresetParam := r.URL.Query().Get("colorPreset")
	colorPaletteRepeat, _ := strconv.Atoi(r.URL.Query().Get("colorPaletteRepeat"))
	colorPaletteLength, _ := strconv.Atoi(r.URL.Query().Get("colorPaletteLength"))
	colorPaletteReverse, _ := strconv.ParseBool(r.URL.Query().Get("colorPaletteReverse"))
	colorPaletteHardStops, _ := strconv.ParseBool(r.URL.Query().Get("colorPaletteHardStops"))

	colorPreset, err := s.colorPresets.GetByIdent(colorPresetParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Unknown color preset")
		return
	}

	var commonFractParams = lib.CommonFractParams{
		ImageWidth:            width,
		ImageHeight:           height,
		CenterCX:              centerCX,
		CenterCY:              centerCY,
		DiameterCX:            diameterCX,
		MaxIterations:         maxIterations,
		ColorPalette:          colorPreset.Palette,
		ColorPaletteRepeat:    colorPaletteRepeat,
		ColorPaletteLength:    colorPaletteLength,
		ColorPaletteReverse:   colorPaletteReverse,
		ColorPaletteHardStops: colorPaletteHardStops,
	}
	juliaKr, _ := strconv.ParseFloat(r.URL.Query().Get("juliaKr"), 64)
	juliaKi, _ := strconv.ParseFloat(r.URL.Query().Get("juliaKi"), 64)
	juliaParams := lib.JuliaFractal{
		JuliaKr: juliaKr,
		JuliaKi: juliaKi,
	}
	streamFractalImage(iterFunc, commonFractParams, juliaParams, format, w)
}

func (s *WebServer) handleWmtsRequest(w http.ResponseWriter, r *http.Request) {
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
	colorPaletteLength, _ := strconv.Atoi(r.URL.Query().Get("colorPaletteLength"))
	if colorPaletteLength == 0 {
		colorPaletteLength = -1
	}

	colorPaletteReverse, _ := strconv.ParseBool(r.URL.Query().Get("colorPaletteReverse"))
	colorPaletteHardStops, _ := strconv.ParseBool(r.URL.Query().Get("colorPaletteHardStops"))

	colorPreset, err := s.colorPresets.GetByIdent(colorPresetParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Unknown color preset")
		return
	}

	var commonFractParams = lib.CommonFractParams{
		ImageWidth:            tileWidthPixels,
		ImageHeight:           tileWidthPixels,
		CenterCX:              centerCX,
		CenterCY:              centerCY,
		DiameterCX:            tileWidthFractal,
		MaxIterations:         maxIterations,
		ColorPalette:          colorPreset.Palette,
		ColorPaletteRepeat:    colorPaletteRepeat,
		ColorPaletteLength:    colorPaletteLength,
		ColorPaletteReverse:   colorPaletteReverse,
		ColorPaletteHardStops: colorPaletteHardStops,
	}
	juliaKr, _ := strconv.ParseFloat(r.URL.Query().Get("juliaKr"), 64)
	juliaKi, _ := strconv.ParseFloat(r.URL.Query().Get("juliaKi"), 64)
	juliaParams := lib.JuliaFractal{
		JuliaKr: juliaKr,
		JuliaKi: juliaKi,
	}
	streamFractalImage(iterFunc, commonFractParams, juliaParams, "jpg", w)
}

// calculates a fractal image and streams it to the given response writer
func streamFractalImage(iterFunc string, commonFractParams lib.CommonFractParams, additionalFractParams any, format string, w http.ResponseWriter) {
	var fractal lib.Fractal

	switch iterFunc {
	case "mandelbrot":
		fractal = lib.NewMandelbrotFractal(commonFractParams)
		break
	case "mandelbrot3":
		fractal = lib.NewMandelbrot3Fractal(commonFractParams)
		break
	case "mandelbrot4":
		fractal = lib.NewMandelbrot4Fractal(commonFractParams)
		break
	case "julia":
		var juliaParams = additionalFractParams.(lib.JuliaFractal)
		fractal = lib.NewJuliaFractal(commonFractParams, juliaParams.JuliaKr, juliaParams.JuliaKi)
		break
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	img := lib.CalcFractalImage(fractal)
	w.Header().Set("Cache-Control", "public, max-age=15552000;")
	switch format {
	case "png":
		w.Header().Set("Content-Type", "image/jpeg")
		img.EncodePng(w)
		break
	case "jpg":
		fallthrough
	case "jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
		img.EncodeJpeg(w)
		break
	}
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

func (s *WebServer) handlePaletteViewer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=15552000;")
	w.WriteHeader(http.StatusOK)
	width, _ := strconv.Atoi(r.URL.Query().Get("width"))
	if width <= 0 {
		width = 1024
	}
	height, _ := strconv.Atoi(r.URL.Query().Get("height"))
	if height <= 0 {
		height = 100
	}
	maxIter, _ := strconv.Atoi(r.URL.Query().Get("maxIterations"))
	if maxIter == 0 {
		maxIter = width
	}
	colorPaletteRepeat, _ := strconv.Atoi(r.URL.Query().Get("paletteRepeat"))
	if colorPaletteRepeat == 0 {
		colorPaletteRepeat = 1
	}
	colorPaletteLength, _ := strconv.Atoi(r.URL.Query().Get("paletteLength"))
	if colorPaletteLength == 0 {
		colorPaletteLength = -1
	}
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		dir = "horizontal"
	}

	switch dir {
	case "horizontal":
		maxIter = width
		break
	case "vertical":
		maxIter = height
		break
	}

	colorPreset, _ := s.colorPresets.GetByIdent(r.URL.Query().Get("colorPreset"))

	img := lib.FractImage{image.NewRGBA(image.Rect(0, 0, width, height))}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			fractParams := lib.CommonFractParams{
				MaxIterations:      maxIter,
				ColorPalette:       colorPreset.Palette,
				ColorPaletteRepeat: colorPaletteRepeat,
				ColorPaletteLength: colorPaletteLength,
			}
			var iterValue float64
			switch dir {
			case "horizontal":
				iterValue = float64(x)
				break
			case "vertical":
				iterValue = float64(y)
				break
			}
			result := lib.FractFunctionResult{}
			lib.SetPaletteColor(&img, x, y, iterValue, fractParams, result)
		}
	}
	img.EncodeJpeg(w)
}

func (s *WebServer) handlePing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}
