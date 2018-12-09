package readtemperature

import (
	"bytes"
	"regexp"
	"strconv"
)

type TemperatureInfo struct {
	Value float32
	Unit string
}

func GetTemp() (TemperatureInfo, error) {
	cmd := ReadTemperatureExecutableCommand()

	result := TemperatureInfo{
		Value:  -1,
		Unit:   "N/A",
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return result, err
	} else {
		re := regexp.MustCompile("temp=(.*)'C")
		output := out.String()
		temperatureAsString := re.FindStringSubmatch(output)
		if len(temperatureAsString) == 2 {
			resultingFloat, err := strconv.ParseFloat(temperatureAsString[1], 32)
			if (err == nil) {
				// YAY!
				result.Value = float32(resultingFloat)
				result.Unit = "Celsius"	// we have celsius here since we were looking for that in our regex.
			}
		}
		return result, nil
	}
}