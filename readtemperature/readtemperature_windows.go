package readtemperature

import "os/exec"

func ReadTemperatureExecutableCommand() *exec.Cmd {
	cmd := exec.Command("C:/TOOLS/raspi_measure_temp_faker.exe")
	cmd.Dir = "C:/TOOLS/"
	return cmd
}