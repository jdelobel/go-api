package media

import (
	"context"
	"fmt"
	"net/url"

	"github.com/go-api/internal/platform/db"
	"github.com/go-api/internal/platform/web"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const mediasTable = "medias"

// List retrieves a list of existing medias from the database.
func List(ctx context.Context, dbConn *db.DB, queryParams url.Values) ([]Media, error) {
	sqlQuery := "SELECT * from medias"
	if len(queryParams) > 0 {
		where := dbConn.BuildWhere(queryParams)
		fmt.Println("where", where)
		sqlQuery += where
	}

	fmt.Println(sqlQuery)
	medias := make([]Media, 0)
	rows, err := dbConn.PSQLQuerier(ctx, sqlQuery)

	if err != nil {
		return nil, errors.Wrap(err, "db.medias.find()")
	}
	for rows.Next() {
		data := Media{}
		err := rows.Scan(
			&data.ID,
			&data.Slug,
			&data.Title,
			&data.URL,
			&data.Specific,
			&data.PublishedAt,
			&data.ExpiredAt,
			&data.Publisher,
			&data.CreatedAt,
			&data.UpdatedAt,
			&data.RestoredAt,
			&data.DeletedAt,
		)

		if err != nil {
			return nil, err
		}
		medias = append(medias, data)
	}
	return medias, nil
}

// Retrieve gets the specified medias from the database.
func Retrieve(ctx context.Context, dbConn *db.DB, mediaID string) (*Media, error) {
	if !bson.IsObjectIdHex(mediaID) {
		return nil, errors.Wrapf(web.ErrInvalidID, "bson.IsObjectIdHex: %s", mediaID)
	}

	q := bson.M{"media_id": mediaID}

	var m *Media
	if _, err := dbConn.PSQLExecute(ctx, "SELECT * from medias where media_id='"+mediaID+"'"); err != nil {
		if err == mgo.ErrNotFound {
			return nil, web.ErrNotFound
		}
		return nil, errors.Wrap(err, fmt.Sprintf("db.medias.find(%s)", db.Query(q)))
	}

	return m, nil
}

// Create inserts a new media into the database.
func Create(ctx context.Context, dbConn *db.DB, cm *CreateMedia) (string, error) {
	params := []interface{}{cm.ID, cm.Title, cm.URL, cm.Slug, cm.Publisher}
	query := "INSERT INTO medias(media_id, title, url, slug, publisher) VALUES($1, $2,$3,$4,$5)"
	if _, err := dbConn.PSQLExecute(ctx, query, params...); err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("db.medias.insert(%s)", db.Query(query)))
	}

	return cm.ID, nil
}

// Update replaces a media document in the database.
func Update(ctx context.Context, dbConn *db.DB, mediaID string, um *CreateMedia) error {
	if !bson.IsObjectIdHex(mediaID) {
		return errors.Wrap(web.ErrInvalidID, "check objectid")
	}

	q := bson.M{"media_id": mediaID}
	m := bson.M{"$set": um}

	if _, err := dbConn.PSQLExecute(ctx, "UPDATE medias"); err != nil {
		if err == mgo.ErrNotFound {
			return web.ErrNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("db.media.update(%s, %s)", db.Query(q), db.Query(m)))
	}

	return nil
}
