package database

import (
	"database/sql"
	"fmt"
	"pitemp/logging"

	/* blank-imported Postgres driver */
	_ "github.com/lib/pq"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

// DBConfig -- Struct for yaml-based DB config
type DBConfig struct {
	Host     string `yaml:"host"`
	Port     uint32 `yaml:"port"`
	DBname   string `yaml:"dbname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	// a unique identifier for distinguishing individual database tables
	TableIdentifier string `yaml:"table-postfix"`
}

// IInserteableMeasurement -- interface to define required methods of inserteable measurements.
type IInserteableMeasurement interface {
	InserteableMeasurementValue() float32
	InserteableMeasurementUnit() string
}

var dbconfig *DBConfig
var mydb *sql.DB
var logger = logging.NewDevLog("database")

// Open -- Opens a database connection according to yaml file 'dbconfig.yml'
func Open(dbconfigArg *DBConfig) {
	var err error
	dbconfig = dbconfigArg
	// fail hard in case of a stupid config
	err = connectDatabase()
	if err != nil {
		panic(err)
	}
	// fail hard in case of a stupid config
	err = ensureTableExists()
	if err != nil {
		panic(err)
	}
}

// Close -- closes the given database connection
func Close() {
	if mydb != nil {
		err := mydb.Close()
		if err != nil {
			logger.Info("DB connection has been shut down gracefully")
		} else {
			logger.Warn("Error closing DB connection")
		}
	}
}

// InsertMeasurement -- insert a measurement
func InsertMeasurement(measurement IInserteableMeasurement) error {
	logger.Debug("Inserting meaurement ...",
		zap.Float32("value", measurement.InserteableMeasurementValue()),
		zap.String("unit", measurement.InserteableMeasurementUnit()))

	tableName := "raspi_measurements_" + dbconfig.TableIdentifier

	statement := "INSERT INTO public." + tableName + " (measurement_timestamp, value, unit) " +
		"VALUES (current_timestamp, $1, $2) RETURNING measurement_id"
	stmt, err := mydb.Prepare(statement)
	defer stmt.Close()
	if err != nil {
		logger.Error("Error preparing statement.", zap.String("statement", statement), zap.Error(err))
		return err
	}
	measurementID := int64(0)
	err = stmt.QueryRow(measurement.InserteableMeasurementValue(), measurement.InserteableMeasurementUnit()).Scan(&measurementID)
	if err != nil {
		logger.Error("Error executing statement.", zap.Error(err))
		return err
	}
	logger.Info("Successfully inserted measurement.",
		zap.Int64("measurement_id", measurementID))

	return nil
}

func connectDatabase() error {
	// assemble CONNECT string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	logger.Info("Connecting to postgres db", zap.String("connection-string", psqlInfo))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	mydb = db

	err = db.Ping()
	if err != nil {
		return err
	}

	logger.Info("Successfully connected to database.")

	return nil
}

/**
tableIdentifier should be the raspi's mac address
*/
func ensureTableExists() error {
	// simple validation
	if len(dbconfig.TableIdentifier) < 1 {
		return errors.New("Cannot use empty table postfix.")
	}

	tableName := "raspi_measurements_" + dbconfig.TableIdentifier
	_, err := mydb.Exec(
		"CREATE TABLE IF NOT EXISTS public." + tableName +
			`(
			measurement_id serial NOT NULL,
			measurement_timestamp timestamp with time zone,
			value numeric NOT NULL,
			unit character varying(255) NOT NULL,
			PRIMARY KEY (measurement_id)
		) WITH (
			OIDS = FALSE
		);
	`)

	if err != nil {
		logger.Error("Error executing CREATE TABLE statement")
		return errors.New("Error executing CREATE TABLE statement")
	}

	logger.Info("Successfully ensured existence of measurement table.", zap.String("tablename", tableName))

	return nil
}
