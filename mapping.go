package nozzle

import (
	"database/sql"
)

// Entity is a class that represent an entry in the DB
type Entity interface {
	Map(row *sql.Row) error
}

// Map maps the DB to an entity
func Map(entity Entity, rows *sql.Row) error {

	return entity.Map(row)
}

// ReturnOne returns one (and only one) entity from the DB
func ReturnOne(entity Entity, sql string, args ...interface{}) error {

	row, err := RunSQL(sql, args...)
	if err != nil {
		logger.Error("SQL Execution", err)
		return err
	}

	return Map(entity, row)
}
