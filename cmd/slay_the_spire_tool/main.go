package main

import (
	"io/fs"
	"log"
	"net/http"
	"strings"

	wailsv2 "github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	embeddedassets "slay_the_spire_tool"
	wailsbootstrap "slay_the_spire_tool/internal/transport/wails"
)

func main() {
	frontendFS, err := fs.Sub(embeddedassets.FS, "frontend/src")
	if err != nil {
		log.Fatalf("prepare embedded frontend fs failed: %v", err)
	}
	assetsFS, err := fs.Sub(embeddedassets.FS, "assets")
	if err != nil {
		log.Fatalf("prepare embedded assets fs failed: %v", err)
	}

	dataFS, err := fs.Sub(embeddedassets.FS, "config/data")
	if err != nil {
		log.Fatalf("prepare embedded data fs failed: %v", err)
	}

	app, err := wailsbootstrap.NewBoundAppWithDataFS(dataFS)
	if err != nil {
		log.Fatalf("init backend app failed: %v", err)
	}

	if err := wailsv2.Run(&options.App{
		Title:  "Slay_the_Spire_2_Tool",
		Width:  1160,
		Height: 760,
		AssetServer: &assetserver.Options{
			Assets: frontendFS,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.HasPrefix(r.URL.Path, "/assets/") {
					http.StripPrefix("/assets/", http.FileServer(http.FS(assetsFS))).ServeHTTP(w, r)
					return
				}
				http.NotFound(w, r)
			}),
		},
		Bind: []interface{}{app},
	}); err != nil {
		log.Fatalf("wails run failed: %v", err)
	}
}
