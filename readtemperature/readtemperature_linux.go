package readtemperature

func ReadTemperatureExecutableFQP() string{
	return "/opt/vc/bin/vcgencmd measure_temp"
}