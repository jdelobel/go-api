package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/url"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	// postgres driver
	_ "github.com/lib/pq"
)

// ErrInvalidDBProvided is returned in the event that an uninitialized db is
// used to perform actions against.
var ErrInvalidDBProvided = errors.New("invalid DB provided")

// DB is a collection of support for different DB technologies. Currently
// only Postgres has been implemented. We want to be able to access the raw
// database support for the given DB so an interface does not work. Each
// database is too different.
type DB struct {
	// Postgres Support.
	database *sqlx.DB
}

// NewPSQL returns a new DB value for use with Postgresql
func NewPSQL(url string) (*DB, error) {
	db, err := newPSQL(url)
	if err != nil {
		return nil, errors.Wrap(err, "NewPSQL")
	}

	return &DB{database: db}, nil
}

// PSQLClose closes a DB value being used with Postgresql.
func (db *DB) PSQLClose() error {
	if err := db.database.Close(); err != nil {
		return errors.Wrapf(err, "PSQLClose: %v", err)
	}
	return nil
}

// PSQLExecute is used to execute Postgres commands.
func (db *DB) PSQLExecute(ctx context.Context, query string, params ...interface{}) (sql.Result, error) {
	if db == nil {
		return nil, errors.Wrap(ErrInvalidDBProvided, "db == nil")
	}
	return db.database.Exec(query, params...)
}

// PSQLQuerier is used to execute Postgres commands.
func (db *DB) PSQLQuerier(ctx context.Context, query string, params ...interface{}) (*sqlx.Rows, error) {
	if db == nil {
		return nil, errors.Wrap(ErrInvalidDBProvided, "db == nil")
	}
	return db.database.Queryx(query, params...)
}

// PSQLQueryRawx is used to execute retrive one raw. Can be used to get raw by its id
func (db *DB) PSQLQueryRawx(ctx context.Context, query string, params ...interface{}) (*sqlx.Row, error) {
	if db == nil {
		return nil, errors.Wrap(ErrInvalidDBProvided, "db == nil")
	}
	return db.database.QueryRowx(query, params...), nil
}

// newPSQL creates a new postgres connection.
func newPSQL(url string) (*sqlx.DB, error) {

	// Create a session which maintains a pool of socket connections
	// to our Postgresql.
	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot connect to database %s", url)
	}

	if err = db.Ping(); err != nil {
		return nil, errors.Wrapf(err, "Cannot ping database %s", url)
	}
	return db, nil
}

// Query provides a string version of the value
func Query(value interface{}) string {
	json, err := json.Marshal(value)
	if err != nil {
		return ""
	}

	return string(json)
}

// Ping checks the connection availability
func (db *DB) Ping() error {
	if db == nil {
		return errors.Wrap(ErrInvalidDBProvided, "db == nil")
	}
	if err := db.database.Ping(); err != nil {
		return errors.Wrap(err, "Failed to ping database")
	}
	return nil
}

// GetQueryOperator identifies operator on a join
func getQueryOperator(op string) (string, error) {
	op = strings.Replace(op, "$", "", -1)
	op = strings.Replace(op, " ", "", -1)

	switch op {
	case "eq":
		return "=", nil
	case "ne":
		return "!=", nil
	case "gt":
		return ">", nil
	case "gte":
		return ">=", nil
	case "lt":
		return "<", nil
	case "lte":
		return "<=", nil
	case "in":
		return "IN", nil
	case "nin":
		return "NOT IN", nil
	case "notnull":
		return "IS NOT NULL", nil
	case "null":
		return "IS NULL", nil
	}
	return "", errors.New("Invalid operator")
}

// BuildWhere builds WHERE statement
func (db *DB) BuildWhere(queryParams url.Values) (string, error) {
	removeOperatorRegex := regexp.MustCompile(`\$[a-z]+.`)
	sqlWhere := " WHERE "
	result := ""
	i := 0
	for k, v := range queryParams {
		value := removeOperatorRegex.ReplaceAllString(v[0], "")
		op, err := getQueryOperator(strings.Split(v[0], ".")[0])
		if err != nil {
			return "", errors.Wrap(err, "BuildWhere")
		}
		if i == 0 {
			result = k + " " + op + "'" + value + "'"
		} else {
			result += " AND " + k + " " + op + "'" + value + "'"
		}
		i++
	}
	sqlWhere += result
	return sqlWhere, nil

}
