package db

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"regexp"
	"strings"

	"database/sql"

	"fmt"
	// postgres driver
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// ErrInvalidDBProvided is returned in the event that an uninitialized db is
// used to perform actions against.
var ErrInvalidDBProvided = errors.New("invalid DB provided")

// MasterDB uggly global postgress conf
var MasterDB DB

// DB is a collection of support for different DB technologies. Currently
// only MongoDB has been implemented. We want to be able to access the raw
// database support for the given DB so an interface does not work. Each
// database is too different.
type DB struct {

	// Postgres Support.
	database *sql.DB
}

// NewPSQL returns a new DB value for use with Postgresql based on a registered
// master session.
func NewPSQL(driver string, url string) (*DB, error) {
	ses, err := newPSQL(driver, url)
	if err != nil {
		return nil, errors.Wrapf(err, "NewPSQL: %s,%s", driver, url)
	}

	db := DB{
		database: ses,
	}

	MasterDB = db

	return &db, nil
}

// PSQLClose closes a DB value being used with Postgresql.
func (db *DB) PSQLClose() {
	db.database.Close()
}

// PSQLExecute is used to execute MongoDB commands.
func (db *DB) PSQLExecute(ctx context.Context, query string, params ...interface{}) (sql.Result, error) {
	if db == nil {
		return nil, errors.Wrap(ErrInvalidDBProvided, "db == nil")
	}
	if len(params) == 0 {
		return db.database.Exec(query)
	}
	return db.database.Exec(query, params...)
}

// PSQLQuerier is used to execute MongoDB commands.
func (db *DB) PSQLQuerier(ctx context.Context, query string, params ...interface{}) (*sql.Rows, error) {
	if db == nil {
		return nil, errors.Wrap(ErrInvalidDBProvided, "db == nil")
	}
	if len(params) == 0 {
		return db.database.Query(query)
	}
	return db.database.Query(query, params...)
}

// Query process queries
func (db *DB) Query(SQL string, params ...interface{}) (jsonData []byte, err error) {
	if db == nil {
		return nil, errors.Wrap(ErrInvalidDBProvided, "db == nil")
	}
	if err != nil {
		log.Println(err)
		return
	}

	SQL = fmt.Sprintf("SELECT json_agg(s) FROM (%s) s", SQL)

	prepare, err := db.database.Prepare(SQL)
	if err != nil {
		return
	}
	defer prepare.Close()

	err = prepare.QueryRow(params...).Scan(&jsonData)

	return
}

// newPSQL creates a new postgres connection.
func newPSQL(driver string, url string) (*sql.DB, error) {

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	ses, err := sql.Open(driver, url)
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot connect to database: %s", url)
	}
	return ses, nil
}

// Query provides a string version of the value
func Query(value interface{}) string {
	json, err := json.Marshal(value)
	if err != nil {
		return ""
	}

	return string(json)
}

// GetQueryOperator identify operator on a join
func getQueryOperator(op string) (string, error) {
	fmt.Println("getQueryOperator", op)
	op = strings.Replace(op, "$", "", -1)
	op = strings.Replace(op, " ", "", -1)
	fmt.Println("getQueryOperator after", op)

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

	err := errors.New("Invalid operator")
	return "", err

}

// BuildWhere build where statement
func (db *DB) BuildWhere(queryParams url.Values) string {
	removeOperatorRegex := regexp.MustCompile(`\$[a-z]+.`)
	sqlWhere := " WHERE "
	result := ""
	i := 0
	for k, v := range queryParams {
		fmt.Println(k, v)
		value := removeOperatorRegex.ReplaceAllString(v[0], "")
		op, err := getQueryOperator(strings.Split(v[0], ".")[0])
		if err != nil {
			fmt.Println("Error", err)
			break
		}
		if i == 0 {
			result = k + " " + op + "'" + value + "'"
		} else {
			result += " AND " + k + " " + op + "'" + value + "'"
		}
		i++
	}
	sqlWhere += result
	return sqlWhere

}
