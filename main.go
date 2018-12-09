package main

import (
	"pitemp/logging"
	"pitemp/readtemperature"
	"sync"
	"time"
)
import "go.uber.org/zap"

var logger = logging.New("pitemp", false)
var allJobsDone = &sync.WaitGroup{}

func main(){
	logger.Info("### STARTUP")

	allJobsDone.Add(1)
	go loopedTemperatureRead("go")

	allJobsDone.Wait()
	logger.Info("### EXIT")
}

func loopedTemperatureRead(foo string) {
	logger.Info("IN:  %s", zap.String("foo", foo))
	defer allJobsDone.Done()

	for {
		info, err := readtemperature.GetTemp()
		if err==nil {
			logger.Info("Current Temperature: ", zap.String("Unit", info.Unit), zap.Float32("Value", info.Value));
		}
		time.Sleep(time.Second)
	}
	logger.Info("OUT: %s", zap.String("foo", foo))
}

