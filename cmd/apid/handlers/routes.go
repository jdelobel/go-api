package handlers

import (
	"net/http"
	"path"
	"runtime"

	"github.com/apex/log"

	"github.com/jdelobel/go-api/config"
	"github.com/jdelobel/go-api/internal/middleware"
	"github.com/jdelobel/go-api/internal/platform/db"
	"github.com/jdelobel/go-api/internal/platform/rabbitmq"
	"github.com/jdelobel/go-api/internal/platform/web"
)

// API returns a handler for a set of routes.
func API(masterDB *db.DB, log *log.Entry, c config.Config, rbmq *rabbitmq.RabbitMQ) http.Handler {

	// Create the web handler for setting routes and middleware.
	app := web.New(log, middleware.RequestLogger, middleware.ErrorHandler)
	// Create the file server to serve static content such as
	// the index.html page.
	statics := http.FileServer(http.Dir(staticsDir()))
	app.TreeMux.NotFoundHandler = statics.ServeHTTP

	// Initialize the routes for the API binding the route to the
	// handler code for each specified verb.
	m := Image{masterDB, rbmq}
	h := Healthzcheck{masterDB}
	s := Swagger{URL: c.AppHost + ":" + c.AppPort}
	app.Handle("GET", "/v1/healthz", h.Healthz)
	app.Handle("GET", "/v1/readiness", h.Readiness)
	app.Handle("GET", "/v1/swagger/swagger.yaml", s.GetAPIDocs)
	app.Handle("GET", "/v1/images", m.List)
	app.Handle("POST", "/v1/images", m.Create)
	app.Handle("GET", "/v1/images/:id", m.Retrieve)
	app.Handle("PUT", "/v1/images/:id", m.Update)
	return app
}

// staticsDir builds a full path to the 'statics' directory
// that is relative to this file. It uses a trick of the
// runtime package to get the path of the file that calls
// this function.
func staticsDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "../statics")
}
