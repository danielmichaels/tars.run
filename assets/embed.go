package assets

import "embed"

//go:embed "static" "migrations" "view"
var Files embed.FS
