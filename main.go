package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"pitemp/database"
	"pitemp/logging"
	"pitemp/readtemperature"
	"syscall"
	"time"

	"github.com/apex/log"
	"go.uber.org/zap"
	yaml "gopkg.in/yaml.v2"
)

var logger = logging.New("pitemp", false)
var temperatureInfoChannel = make(chan readtemperature.TemperatureInfo)

func main() {
	logger.Info("### STARTUP")

	// INIT
	var cfg database.DBConfig
	readDatabaseConfig(&cfg)
	database.Open(&cfg)
	defer database.Close()

	// GO
	go loopedTemperatureRead()
	go temperaturePrinter()

	// wait indefinitely until external abortion
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) // Ctrl + c
	<-sigs
	logger.Info("### EXIT")
}

// ==== I/O and properties ====

// read config from 'dbconfig.yml'
func readDatabaseConfig(dbconfig *database.DBConfig) {
	var err error
	var bytes []byte
	bytes, err = ioutil.ReadFile("dbconfig.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(bytes, dbconfig)
	if err != nil {
		panic(err)
	}
	log.Info("DBConfig parsed.")
}

func temperaturePrinter() {
	for {
		info := <-temperatureInfoChannel

		logger.Info("Current Temperature: ", zap.String("Unit", info.Unit), zap.Float32("Value", info.Value))

		err := database.InsertMeasurement(info)
		if err != nil {
			logger.Error("Cannot persist measurement")
		}

	}
}

func temperatureRead() {
	info, err := readtemperature.GetTemp()
	if err == nil {
		//logger.Info("Current Temperature: ", zap.String("Unit", info.Unit), zap.Float32("Value", info.Value));
		temperatureInfoChannel <- info
	} else {
		logger.Error("Could not read temperature: ", zap.Error(err))
	}

}

func loopedTemperatureRead() {
	for {
		select {
		case <-time.After(5 * time.Second):
			temperatureRead()
		}
	}
}
