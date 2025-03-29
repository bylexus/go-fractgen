package web

import (
	"net/http"
	"runtime"
	"strconv"

	"github.com/bylexus/go-fract/lib"
	"github.com/bylexus/go-stdlib/ethreads"
)

type WebServer struct {
	http.Server

	threadPool *ethreads.ThreadPool

	colorPresets   []lib.ColorPreset
	fractalPresets []lib.FractalPreset
}

func NewWebServer(colorPresets []lib.ColorPreset, fractalPresets []lib.FractalPreset) *WebServer {
	threadPool := ethreads.NewThreadPool(runtime.NumCPU(), nil)
	threadPool.Start()
	server := &WebServer{
		threadPool:     &threadPool,
		colorPresets:   colorPresets,
		fractalPresets: fractalPresets,
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("webroot")))
	mux.HandleFunc("/fractal-image.png", server.handleFractalImage)

	listenAddr := ":8000"
	var handler http.Handler = NewLogHandler(mux)
	server.Server = http.Server{
		Addr: listenAddr, Handler: handler,
	}

	return server
}

func (s *WebServer) handleFractalImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	width, _ := strconv.Atoi(r.URL.Query().Get("width"))
	height, _ := strconv.Atoi(r.URL.Query().Get("height"))
	maxIterations, _ := strconv.Atoi(r.URL.Query().Get("maxIterations"))
	// iterFunc := strings.ToLower(r.URL.Query().Get("iterFunc"))
	centerCX, _ := strconv.ParseFloat(r.URL.Query().Get("centerCX"), 64)
	centerCY, _ := strconv.ParseFloat(r.URL.Query().Get("centerCY"), 64)
	diameterCX, _ := strconv.ParseFloat(r.URL.Query().Get("diameterCX"), 64)

	fractal := lib.NewMandelbrotFractal(
		width, height,
		centerCX, centerCY, diameterCX,
		maxIterations, s.colorPresets[0].Palette,
	)

	img := fractal.CalcFractalImage(s.threadPool)
	w.WriteHeader(http.StatusOK)
	img.EncodePng(w)
}
