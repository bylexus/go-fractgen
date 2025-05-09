package web

import (
	"encoding/json"
	"fmt"
	"image"
	"math/big"
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
	iterFunc, _ := lib.FractalTypeFromString(r.URL.Query().Get("iterFunc"))
	centerCX, _, _ := big.ParseFloat(r.URL.Query().Get("centerCX"), 10, lib.SYS_PRECISION, big.ToNearestEven)
	centerCY, _, _ := big.ParseFloat(r.URL.Query().Get("centerCY"), 10, lib.SYS_PRECISION, big.ToNearestEven)
	diameterCX, _, _ := big.ParseFloat(r.URL.Query().Get("diameterCX"), 10, lib.SYS_PRECISION, big.ToNearestEven)
	colorPresetParam := r.URL.Query().Get("colorPreset")
	colorPaletteRepeat, _ := strconv.Atoi(r.URL.Query().Get("colorPaletteRepeat"))
	colorPaletteLength, _ := strconv.Atoi(r.URL.Query().Get("colorPaletteLength"))
	colorPaletteReverse, _ := strconv.ParseBool(r.URL.Query().Get("colorPaletteReverse"))
	juliaKr, _, _ := big.ParseFloat(r.URL.Query().Get("juliaKr"), 10, lib.SYS_PRECISION, big.ToNearestEven)
	juliaKi, _, _ := big.ParseFloat(r.URL.Query().Get("juliaKi"), 10, lib.SYS_PRECISION, big.ToNearestEven)

	colorPreset, err := s.colorPresets.GetByIdent(colorPresetParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Unknown color preset")
		return
	}

	var fractal lib.Fractal
	var commonFractParams = lib.CommonFractParams{
		ImageWidth:          width,
		ImageHeight:         height,
		CenterCX:            centerCX,
		CenterCY:            centerCY,
		DiameterCX:          diameterCX,
		MaxIterations:       maxIterations,
		ColorPalette:        colorPreset.Palette,
		ColorPaletteRepeat:  colorPaletteRepeat,
		ColorPaletteLength:  colorPaletteLength,
		ColorPaletteReverse: colorPaletteReverse,
	}

	fractal, err = lib.NewFractalFromParams(iterFunc, commonFractParams, juliaKr, juliaKi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	streamFractalImage(w, fractal, format)
}

func (s *WebServer) handleWmtsRequest(w http.ResponseWriter, r *http.Request) {
	zoomLevel, _ := strconv.Atoi(r.URL.Query().Get("TileMatrix"))
	if zoomLevel < 0 || zoomLevel > 50 {
		zoomLevel = 0
	}

	iterFunc, err := lib.FractalTypeFromString(r.URL.Query().Get("iterFunc"))

	tileX, _ := strconv.Atoi(r.URL.Query().Get("TileCol"))
	tileY, _ := strconv.Atoi(r.URL.Query().Get("TileRow"))

	originX := big.NewFloat(-1.7).SetPrec(lib.SYS_PRECISION)
	originY := big.NewFloat(-1.0).SetPrec(lib.SYS_PRECISION)

	tileWidthPixels, _ := strconv.Atoi(r.URL.Query().Get("tileWidthPixels"))
	if tileWidthPixels == 0 {
		tileWidthPixels = 256
	}
	resolution, _, _ := big.ParseFloat(r.URL.Query().Get("resolution"), 10, lib.SYS_PRECISION, big.ToNearestEven)

	tileWidthFractal := new(big.Float).SetPrec(lib.SYS_PRECISION)
	tileWidthFractal.Mul(resolution, big.NewFloat(float64(tileWidthPixels)))
	halfTileWidthFractal := new(big.Float).SetPrec(lib.SYS_PRECISION)
	halfTileWidthFractal.Quo(tileWidthFractal, big.NewFloat(2.0))

	// centerCX := originX + float64(tileX)*tileWidthFractal + (tileWidthFractal / 2)
	tileXStart := new(big.Float).SetPrec(lib.SYS_PRECISION).Copy(tileWidthFractal)
	tileXStart.Mul(tileXStart, big.NewFloat(float64(tileX)))
	tileXStart.Add(tileXStart, originX)
	centerCX := tileXStart.Add(tileXStart, halfTileWidthFractal)

	// centerCY := originY + float64(-1*tileY)*tileWidthFractal - (tileWidthFractal / 2)
	tileYStart := new(big.Float).SetPrec(lib.SYS_PRECISION).Copy(tileWidthFractal)
	tileYStart.Mul(tileYStart, big.NewFloat(float64(-1*tileY)))
	tileYStart.Add(tileYStart, originY)
	centerCY := tileYStart.Sub(tileYStart, halfTileWidthFractal)

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

	colorPreset, err := s.colorPresets.GetByIdent(colorPresetParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Unknown color preset")
		return
	}
	juliaKr, _, _ := big.ParseFloat(r.URL.Query().Get("juliaKr"), 10, lib.SYS_PRECISION, big.ToNearestEven)
	juliaKi, _, _ := big.ParseFloat(r.URL.Query().Get("juliaKi"), 10, lib.SYS_PRECISION, big.ToNearestEven)

	var fractal lib.Fractal
	var commonFractParams = lib.CommonFractParams{
		ImageWidth:          tileWidthPixels,
		ImageHeight:         tileWidthPixels,
		CenterCX:            centerCX,
		CenterCY:            centerCY,
		DiameterCX:          tileWidthFractal,
		MaxIterations:       maxIterations,
		ColorPalette:        colorPreset.Palette,
		ColorPaletteRepeat:  colorPaletteRepeat,
		ColorPaletteLength:  colorPaletteLength,
		ColorPaletteReverse: colorPaletteReverse,
	}
	fractal, err = lib.NewFractalFromParams(iterFunc, commonFractParams, juliaKr, juliaKi)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	streamFractalImage(w, fractal, "jpg")
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

func streamFractalImage(w http.ResponseWriter, fractal lib.Fractal, format string) {
	img := lib.CalcFractalImage(fractal)
	w.Header().Set("Cache-Control", "public, max-age=15552000;")
	w.WriteHeader(http.StatusOK)
	switch format {
	case "png":
		w.Header().Set("Content-Type", "image/png")
		img.EncodePng(w)
		break
	case "jpeg":
		fallthrough
	case "jpg":
		w.Header().Set("Content-Type", "image/jpeg")
		img.EncodeJpeg(w)
		break
	}
}
