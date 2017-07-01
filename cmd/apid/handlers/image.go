package handlers

import (
	"context"
	"net/http"

	"fmt"

	"github.com/jdelobel/go-api/internal/image"
	"github.com/jdelobel/go-api/internal/platform/db"
	"github.com/jdelobel/go-api/internal/platform/rabbitmq"
	"github.com/jdelobel/go-api/internal/platform/web"
	"github.com/pkg/errors"
)

// Image represents the Image API method handler set.
type Image struct {
	MasterDB *db.DB
	rbmq     *rabbitmq.RabbitMQ

	// ADD OTHER STATE LIKE THE LOGGER AND CONFIG HERE.
}

// List returns all the existing images in the system.
// 200 Success, 404 Not Found, 500 Internal
func (m *Image) List(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	reqDB := m.MasterDB
	qp := r.URL.Query()
	images, err := image.List(ctx, reqDB, qp)
	if err != nil {
		return errors.Wrap(err, "")
	}
	web.Respond(ctx, w, images, http.StatusOK)
	return nil
}

// Retrieve returns the specified image from the system.
// 200 Success, 400 Bad Request, 404 Not Found, 500 Internal
func (m *Image) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	reqDB := m.MasterDB
	image, err := image.Retrieve(ctx, reqDB, params["id"])
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			err = fmt.Errorf("No image found with id %s", params["id"])
			web.RespondError(ctx, w, err, http.StatusNotFound)
			return nil
		}

		return errors.Wrapf(err, "Id: %s", params["id"])
	}

	web.Respond(ctx, w, image, http.StatusOK)
	return nil
}

// Create inserts a new image into the system.
// 200 OK, 400 Bad Request, 500 Internal
func (m *Image) Create(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	reqDB := m.MasterDB
	rbmq := m.rbmq
	var med image.CreateImage
	if err := web.Unmarshal(r.Body, &med); err != nil {
		return errors.Wrap(err, "")
	}

	img, err := image.Create(ctx, reqDB, rbmq, &med)
	if err != nil {
		return errors.Wrapf(err, "Image: %+v", &med)
	}

	web.Respond(ctx, w, img, http.StatusCreated)
	return nil
}

// Update updates the specified image in the system.
// 200 Success, 400 Bad Request, 500 Internal
func (m *Image) Update(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	reqDB := m.MasterDB

	var med image.CreateImage
	if err := web.Unmarshal(r.Body, &med); err != nil {
		return errors.Wrap(err, "")
	}

	if err := image.Update(ctx, reqDB, params["id"], &med); err != nil {
		return errors.Wrapf(err, "Id: %s  Image: %+v", params["id"], &med)
	}

	web.Respond(ctx, w, nil, http.StatusNoContent)
	return nil
}
