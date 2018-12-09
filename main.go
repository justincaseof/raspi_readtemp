package main

import (
	"log"
	"pitemp/readtemperature"
)
import "go.uber.org/zap"

func main(){
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	info, _ := readtemperature.GetTemp()

	logger.Info("Current Temperature: ", zap.String("Unit", info.Unit), zap.Float32("Value", info.Value) );


}

