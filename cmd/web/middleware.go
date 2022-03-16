package main

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// If there's a panic, set "Connection: close" on the response.
				// This will tell Go's HTTP server to automatically close the
				// current connection after the response has been sent.
				w.Header().Set("Connection", "close")
				// The value returned by recover() has a type interface{}, so we
				// use fmt.Errorf() to normalize it into an error and call our
				// serverErrorResponse() helper. This will log the error using
				// our custom Logger type at the ERROR level and send the client
				// a 500 status.
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

// Open will not render the filesystem when navigating to url's
// such as /static/css, instead it will return a 404.
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

// addDefaultData is a helper which will pre-fill the templateData struct with
// default information that is used across several templates.
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.AppName = app.config.Server.AppName
	return td
}

// render is a template rendering helper. It uses a template cache to speed up delivery of templates
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.template[name]
	if !ok {
		http.Error(w, "Template does not exist", 500)
		return
	}
	buf := new(bytes.Buffer)
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	buf.WriteTo(w)
}
