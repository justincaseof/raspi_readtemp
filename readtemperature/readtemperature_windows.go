package readtemperature

import "os/exec"

func ReadTemperatureExecutableCommand() *exec.Cmd {
	return exec.Command("C:/TOOLS/raspi_measure_temp_faker.exe")
}