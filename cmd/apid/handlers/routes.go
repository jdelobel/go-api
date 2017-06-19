package handlers

import (
	"net/http"
	"path"
	"runtime"

	"github.com/go-api/internal/middleware"
	"github.com/go-api/internal/platform/db"
	"github.com/go-api/internal/platform/web"
)

// API returns a handler for a set of routes.
func API(masterDB *db.DB) http.Handler {

	// Create the web handler for setting routes and middleware.
	app := web.New(middleware.RequestLogger, middleware.ErrorHandler)

	// Create the file server to serve static content such as
	// the index.html page.
	views := http.FileServer(http.Dir(viewsDir()))
	app.TreeMux.NotFoundHandler = views.ServeHTTP

	// Initialize the routes for the API binding the route to the
	// handler code for each specified verb.
	m := Media{masterDB}
	app.Handle("GET", "/v1/medias", m.List)
	app.Handle("POST", "/v1/medias", m.Create)
	app.Handle("GET", "/v1/medias/:id", m.Retrieve)
	app.Handle("PUT", "/v1/medias/:id", m.Update)

	return app
}

// viewsDir builds a full path to the 'views' directory
// that is relative to this file. It uses a trick of the
// runtime package to get the path of the file that calls
// this function.
func viewsDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "../views")
}
