package image

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"

	"github.com/apex/log"
	"github.com/jdelobel/go-api/internal/platform/db"
	"github.com/jdelobel/go-api/internal/platform/rabbitmq"
	"github.com/jdelobel/go-api/internal/platform/web"
	"github.com/pkg/errors"
)

// List retrieves a list of existing images from the database.
func List(ctx context.Context, dbConn *db.DB, queryParams url.Values) ([]Image, error) {
	sqlQuery := "SELECT * from images"
	if len(queryParams) > 0 {
		where, err := dbConn.BuildWhere(queryParams)
		if err != nil {
			return nil, errors.Wrap(err, "List")
		}
		sqlQuery += where
	}

	images := make([]Image, 0)
	rows, err := dbConn.PSQLQuerier(ctx, sqlQuery)

	if err != nil {
		return nil, errors.Wrap(err, "List")
	}
	for rows.Next() {
		data := Image{}
		err := rows.StructScan(&data)

		if err != nil {
			return nil, err
		}
		images = append(images, data)
	}
	return images, nil
}

// Retrieve gets the specified images from the database.
func Retrieve(ctx context.Context, dbConn *db.DB, imageID string) (*Image, error) {
	if !IsValidUUID(imageID) {
		return nil, errors.Wrapf(web.ErrInvalidID, "IsValidUUID: %s", imageID)
	}
	row, err := dbConn.PSQLQueryRawx(ctx, "SELECT * from images where id=$1", imageID)

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.images.find(%s)", db.Query(imageID)))
	}
	image := Image{}
	err = row.StructScan(&image)

	if err != nil {
		return nil, err
	}
	return &image, nil
}

// Create inserts a new image into the database.
func Create(ctx context.Context, dbConn *db.DB, rbmq *rabbitmq.RabbitMQ, cm *CreateImage) (*Image, error) {
	params := []interface{}{cm.Title, cm.URL, cm.Slug, cm.Publisher}
	query := "INSERT INTO images(title, url, slug, publisher) VALUES($1,$2,$3,$4) RETURNING *"
	row, err := dbConn.PSQLQueryRawx(ctx, query, params...)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.images.insert(%s)", db.Query(query)))
	}
	var img Image
	if err = row.StructScan(&img); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.images.insert(%s)StructScan", db.Query(query)))
	}
	qn := "image_created"
	imgJSON, err := json.Marshal(&img)
	if err != nil {
		log.Warnf("RabbitMQ: failed to marshal obj %v for image creation: %v", &img, err)
	}
	rbmq.DeclareQueue(qn)
	if err != nil {
		log.Warnf("RabbitMQ: cannot declare queue for image creation: %v", err)
	}
	err = rbmq.Publish(&qn, imgJSON)
	if err != nil {
		log.Warnf("RabbitMQ: failed to publish a message for image creation: %v", err)
	}
	return &img, nil
}

// Update replaces an image document in the database.
func Update(ctx context.Context, dbConn *db.DB, imageID string, um *CreateImage) error {
	if !IsValidUUID(imageID) {
		return errors.Wrapf(web.ErrInvalidID, "IsValidUUID: %s", imageID)
	}

	if _, err := dbConn.PSQLExecute(ctx, "UPDATE images"); err != nil {
		return errors.Wrap(err, fmt.Sprintf("db.image.update(%s, %s)", db.Query(imageID), db.Query(um)))
	}

	return nil
}

// IsValidUUID check if uuid is in valid format
func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
