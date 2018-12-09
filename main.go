package main

import (
	"pitemp/logging"
	"pitemp/readtemperature"
	"time"
)
import "go.uber.org/zap"

var logger = logging.New("pitemp", false)

func main(){
	logger.Info("### STARTUP")
	loopedTemperatureRead()
}

func loopedTemperatureRead() {
	for {
		info, err := readtemperature.GetTemp()
		if err==nil {
			logger.Info("Current Temperature: ", zap.String("Unit", info.Unit), zap.Float32("Value", info.Value));
		}
		time.Sleep(time.Second)
	}
}

