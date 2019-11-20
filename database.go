package nozzle

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/caarlos0/env"
	"github.com/google/uuid"

	// Used internally bu database/sql
	"github.com/lib/pq"
)

const (
	// DuplicateError record in DB
	DuplicateError = 409

	// InternalError error (Default)
	InternalError = 500
)

const (
	local  = "local"
	cloud  = "cloud"
	driver = "postgres"
)

type DatabaseInfo struct {
	Mode             string `env:"DB_MODE"`
	Type             string `env:"DB_TYPE"`
	ConnectionString string `env:"DB_CONNECTIONSTRING"`
}

func (info *databaseInfo) getPassword() string {

	return os.Getenv("DB_PASSWORD")
}

var config = loadConfig()

func loadConfig() databaseInfo {

	config := databaseInfo{}

	err := env.Parse(&config)
	if err != nil {
		logger.Error("Loading Config", err)
		panic(err)
	}

	return config
}

// CreateID returns a UUID
func CreateID() uuid.UUID {

	id, err := uuid.NewRandom()
	if err != nil {
		logger.Error("Creating ID", err)
		panic(err)
	}

	return id
}

func openDB() (*sql.DB, string, error) {

	var connectionString string

	if config.Mode == local {
		connectionString = fmt.Sprintf(config.ConnectionString, config.getPassword())
	} else {
		connectionString = config.ConnectionString

	}

	logger.Info(LogEntry{Action: "Attempting DB Open", Message: connectionString})
	db, err := sql.Open(config.Type, connectionString)
	return db, connectionString, err
}

// RunSQL is used to add an entity to the database
func RunSQL(sql string, args ...interface{}) (*sql.Row, error) {

	db, connectionString, err := openDB()
	if err != nil {
		logger.Error("DB Open Failed", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Error("DB Ping Failed", err)
		return nil, err
	}
	logger.Info(LogEntry{Action: "DB Ping Succeeded", Message: connectionString})

	row := db.QueryRow(sql, args...)
	logger.Info(LogEntry{Action: "SQL Run", Message: sql})

	defer db.Close()

	return row, err
}

func lookupDBError(err error) (int, string) {

	pqErr := err.(*pq.Error)

	switch pqErr.Code {
	case "23505":
		logger.Error("Mapping Error To Duplicate Error", err)
		return DuplicateError, "Duplicate"
	default:
		logger.Error("Mapping Error To Internal Error", err)
		return InternalError, "Internal Error"
	}
}
