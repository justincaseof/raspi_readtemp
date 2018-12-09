package main

import (
	"os"
	"os/signal"
	"pitemp/logging"
	"pitemp/readtemperature"
	"syscall"
	"time"
)
import "go.uber.org/zap"

var logger = logging.New("pitemp", false)
var temperatureInfoChannel = make(chan readtemperature.TemperatureInfo)

func main(){
	logger.Info("### STARTUP")

	go loopedTemperatureRead()
	go temperaturePrinter()

	// wait indefinately until external abortion
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)	// Ctrl + c
	<-sigs
	logger.Info("### EXIT")
}

func temperaturePrinter() {
	for {
		info := <- temperatureInfoChannel
		logger.Info("Current Temperature: ", zap.String("Unit", info.Unit), zap.Float32("Value", info.Value));
	}
}

func temperatureRead() {
	info, err := readtemperature.GetTemp()
	if err==nil {
		//logger.Info("Current Temperature: ", zap.String("Unit", info.Unit), zap.Float32("Value", info.Value));
		temperatureInfoChannel <- info
	} else {
		logger.Error("Could not read temperature: ", zap.Error(err))
	}

}

func loopedTemperatureRead() {
	for {
		select {
		case <-time.After(1 * time.Second):
			temperatureRead()
		}
	}
}