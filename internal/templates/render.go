package templates

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
)

func Render(ctx context.Context, w http.ResponseWriter, status int, t templ.Component) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/html")
	return t.Render(ctx, w)
}
