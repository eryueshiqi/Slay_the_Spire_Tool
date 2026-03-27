package embeddedassets

import "embed"

// FS contains embedded frontend static assets and default JSON data.
//
//go:embed all:frontend/src all:config/data all:assets
var FS embed.FS
