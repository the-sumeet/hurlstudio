package main

import (
	"embed"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

type FileLoader struct {
	http.Handler
	app *App
}

func NewFileLoader(app *App) *FileLoader {
	return &FileLoader{app: app}
}

func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	requestedFilename := strings.TrimPrefix(req.URL.Path, "/")
	println("Requesting file:", requestedFilename)

	// URL decode the filename to handle spaces and special characters
	decodedFilename, err := url.PathUnescape(requestedFilename)
	if err != nil {
		println("Error decoding path:", err.Error())
		decodedFilename = requestedFilename
	}

	// Resolve the path
	var resolvedPath string

	// Check if it's an absolute path using filepath.IsAbs
	// We need to check the original URL path (with leading /) for proper absolute path detection
	originalPath := req.URL.Path
	decodedOriginalPath, err := url.PathUnescape(originalPath)
	if err != nil {
		decodedOriginalPath = originalPath
	}

	if filepath.IsAbs(decodedOriginalPath) {
		// It's an absolute path, use it directly
		resolvedPath = decodedOriginalPath
	} else {
		// Relative path - resolve based on current file or directory
		var baseDir string
		if h.app.currentFile != "" {
			// Use current file's directory as base
			baseDir = filepath.Dir(h.app.currentFile)
		} else if h.app.currentDir != "" {
			// Use current directory as base
			baseDir = h.app.currentDir
		}

		if baseDir != "" {
			resolvedPath = filepath.Join(baseDir, decodedFilename)
		} else {
			resolvedPath = decodedFilename
		}
	}

	resolvedPath = filepath.Clean(resolvedPath)
	println("Resolved path:", resolvedPath)

	fileData, err := os.ReadFile(resolvedPath)
	if err != nil {
		println("Error reading file:", err.Error())
		res.WriteHeader(http.StatusNotFound)
		res.Write([]byte(fmt.Sprintf("Could not load file %s: %v", resolvedPath, err)))
		return
	}

	res.Write(fileData)
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "hurlstudio",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: NewFileLoader(app),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
