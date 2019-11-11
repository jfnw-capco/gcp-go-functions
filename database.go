package gcpfunctions

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/caarlos0/env"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	// DuplicateError record in DB
	DuplicateError = 409

	// InternalError error (Default)
	InternalError = 500
)

const (
	driver   = "postgres"
	host     = "host"
	port     = "port"
	user     = "user"
	password = "password"
	dbName   = "dbname"
	sslMode  = "sslmode"
	mask     = "*"
)

type databaseInfo struct {
	IsCloudSQL bool   `env:"DB_IS_CLOUD_SQL"`
	Type       string `env:"DB_TYPE"`
	Host       string `env:"DB_HOST"`
	Instance   string `env:"DB_INSTANCE"`
	Name       string `env:"DB_NAME"`
	Port       string `env:"DB_PORT"`
	User       string `env:"DB_USER"`
	SSLMode    string `env:"DB_SSLMODE"`
}

var config = loadConfig()

func loadConfig() databaseInfo {

	config := databaseInfo{}

	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}

	return config
}

func (info *databaseInfo) getPassword() string {

	return os.Getenv("DB_PASSWORD")
}

// CreateID returns a UUID
func CreateID() uuid.UUID {

	id, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	return id
}

func openDB() (*sql.DB, error) {

	var connectionString string

	if config.IsCloudSQL {
		connectionString = fmt.Sprintf("host=/cloudsql/%s user=%s password=%s dbname=%s sslmode=%s",
			config.Instance,
			config.User,
			config.getPassword(),
			config.Name,
			config.SSLMode)
	} else {
		connectionString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Host,
			config.Port,
			config.User,
			config.getPassword(),
			config.Name,
			config.SSLMode)
	}

	logger.DebugMessage("DB ConnectionString", connectionString)

	return sql.Open(config.Type, connectionString)
}

// RunSQL is used to add an entity to the database
func RunSQL(sql string, args ...interface{}) (*sql.Row, error) {

	db, err := openDB()
	if err != nil {
		return nil, err
	}

	row := db.QueryRow(sql, args...)
	defer db.Close()

	return row, err
}

func lookupDBError(err error) (int, string) {

	pqErr := err.(*pq.Error)

	switch pqErr.Code {
	case "23505":
		return DuplicateError, "Duplicate"
	default:
		return InternalError, "Internal Error"
	}
}
