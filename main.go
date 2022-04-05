package nmcli

import (
	"errors"
	"fmt"
	"os/exec"
)

func Run() {
	fmt.Println("Le package")
}

// Checks whether nmcli is installed on the system
// Return an error if nmcli not installed, and a version string if it is installed
func ValidateNmcliInstalled() (string, error) {
	// test for precence of nmcli
	msg, err := exec.Command(
		"bash",
		"-c",
		"nmcli --version",
	).CombinedOutput()
	if err != nil {

		return "", errors.New("nmcli not found on this system")
	}

	return string(msg), nil
}
