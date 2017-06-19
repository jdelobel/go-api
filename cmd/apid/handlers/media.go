package handlers

import (
	"context"
	"net/http"

	"github.com/go-api/internal/media"
	"github.com/go-api/internal/platform/db"
	"github.com/go-api/internal/platform/web"
	"github.com/pkg/errors"
)

// Media represents the Media API method handler set.
type Media struct {
	MasterDB *db.DB

	// ADD OTHER STATE LIKE THE LOGGER AND CONFIG HERE.
}

// List returns all the existing medias in the system.
// 200 Success, 404 Not Found, 500 Internal
func (m *Media) List(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	reqDB := m.MasterDB
	qp := r.URL.Query()
	medias, err := media.List(ctx, reqDB, qp)
	if err != nil {
		return errors.Wrap(err, "")
	}
	//reqDB.PSQLClose()
	web.Respond(ctx, w, medias, http.StatusOK)
	return nil
}

// Retrieve returns the specified media from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (m *Media) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	reqDB := m.MasterDB
	//defer reqDB.PSQLClose()

	media, err := media.Retrieve(ctx, reqDB, params["id"])
	if err != nil {
		return errors.Wrapf(err, "Id: %s", params["id"])
	}

	web.Respond(ctx, w, media, http.StatusOK)
	return nil
}

// Create inserts a new media into the system.
// 200 OK, 400 Bad Request, 500 Internal
func (m *Media) Create(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	reqDB := m.MasterDB

	var med media.CreateMedia
	if err := web.Unmarshal(r.Body, &med); err != nil {
		return errors.Wrap(err, "")
	}

	nUsr, err := media.Create(ctx, reqDB, &med)
	if err != nil {
		return errors.Wrapf(err, "Media: %+v", &med)
	}
	//reqDB.PSQLClose()

	web.Respond(ctx, w, nUsr, http.StatusCreated)
	return nil
}

// Update updates the specified media in the system.
// 200 Success, 400 Bad Request, 500 Internal
func (m *Media) Update(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	reqDB := m.MasterDB
	//defer reqDB.PSQLClose()

	var med media.CreateMedia
	if err := web.Unmarshal(r.Body, &med); err != nil {
		return errors.Wrap(err, "")
	}

	if err := media.Update(ctx, reqDB, params["id"], &med); err != nil {
		return errors.Wrapf(err, "Id: %s  Media: %+v", params["id"], &med)
	}

	web.Respond(ctx, w, nil, http.StatusNoContent)
	return nil
}
