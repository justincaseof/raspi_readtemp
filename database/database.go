package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"pitemp/logging"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     uint32 `yaml:"port"`
	DBname   string `yaml:"dbname"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type IInserteableMeasurement interface {
	InserteableMeasurementValue() float32
	InserteableMeasurementUnit()  string
}

var tableIdentifier string
var dbconfig DBConfig
var mydb *sql.DB
var log = logging.NewDevLog("database")

func Open(tableIdentifierArg string) {
	readConfig(&dbconfig)
	initDatabase(tableIdentifierArg)
	ensureTableExists()
}

func Close() {
	if mydb != nil {
		mydb.Close()
	}
}

func InsertMeasurement(measurement IInserteableMeasurement) (error) {
	log.Info("Inserting meaurement ...",
		zap.Float32("value", measurement.InserteableMeasurementValue()),
		zap.String("unit", measurement.InserteableMeasurementUnit()) )



	return errors.New("Could not insert measurement")
}

func initDatabase(tableIdentifierArg string) {
	// store tableIdentifier. should be instance-unique in order to keep measurements of different devices apart.
	tableIdentifier = tableIdentifierArg

	// assemble CONNECT string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	mydb = db

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Info("Successfully connected!")
}

func readConfig(dbconfig *DBConfig) {
	var err error
	var bytes []byte
	bytes, err = ioutil.ReadFile("dbconfig.yml")
	if (err != nil) {
		panic(err)
	}
	err = yaml.Unmarshal(bytes, dbconfig)
	if (err != nil) {
		panic(err)
	}
	log.Info("DBConfig parsed.")
}

/**
 tableIdentifier should be the raspi's mac address
 */
func ensureTableExists() {
	tableName := "raspi_measurements_" + tableIdentifier
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
		log.Error("Error executing CREATE TABLE statement")
	}
}
