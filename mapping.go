package nozzle

import (
	"database/sql"
)

// Entity is a class that represent an entry in the DB
type Entity interface {
	Map(rows *sql.Rows) error
}

// Map maps the DB to an entity
func Map(entity Entity, rows *sql.Rows) error {

	return entity.Map(rows)
}

// ReturnOne returns one (and only one) entity from the DB
func ReturnOne(entity Entity, sql string, args ...interface{}) error {

	rows, err := RunSQL(sql, args...)
	if err != nil {
		logger.Error("SQL Execution", err)
		return err
	}

	return Map(entity, rows)
}
