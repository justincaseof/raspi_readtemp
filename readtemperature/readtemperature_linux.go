package readtemperature

import "os/exec"

func ReadTemperatureExecutableCommand() *exec.Cmd {
	return exec.Command("/opt/vc/bin/vcgencmd", "measure_temp")
}