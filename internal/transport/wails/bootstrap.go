package wails

import (
	"io/fs"
	"slay_the_spire_tool/internal/app"
)

const DefaultDataPath = "config/data"

func NewBoundApp() (*app.App, error) {
	return app.New(DefaultDataPath)
}

func NewBoundAppWithDataFS(dataFS fs.FS) (*app.App, error) {
	return app.NewWithDataFS(dataFS)
}
