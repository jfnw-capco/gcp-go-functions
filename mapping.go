package gcpfunctions

import (
	"database/sql"
	"reflect"
)

// Entity is a class that represent an entry in the DB
type Entity interface {
	Map(row *sql.Row) error
}

// Map maps the DB to an entity
func Map(entity Entity, row *sql.Row) error {

	return entity.Map(row)
}

// ReturnOne returns one (and only one) entity from the DB
func ReturnOne(entity Entity, sql string, args ...interface{}) error {

	row, err := RunSQL(sql, args...)
	if err != nil {
		return err
	}

	return Map(entity, row)
}

// ReturnMany returns one (and only one) entity from the DB
func ReturnMany(entity Entity, sql string, args ...interface{}) error {

	for i := 0; i < 10; i++ {
		reflect.New(reflect.TypeOf(entity))
	}

	return nil
}
