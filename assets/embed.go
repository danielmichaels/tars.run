package assets

import "embed"

//go:embed "html" "static" "migrations"
var Files embed.FS
