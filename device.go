package nmcli

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type Device struct {
	device      string
	device_type string
	state       string
	conn        string
}

func Devices() []Device {
	res, err := exec.Command("nmcli", "device").Output()
	if err != nil {
		// handle error
		fmt.Println(err)
	}
	// process result
	results := make([]Device, 0)
	input := strings.Split(strings.TrimSpace(string(res[:])), "\n")
	// fmt.Printf("%+v\n", input)
	// pop first row (headers)
	for _, line := range input[1:] {
		// fmt.Println(line)
		results = append(results, parseDevice(line))
	}

	return results
}

func parseDevice(dev_line string) Device {
	regex := regexp.MustCompile(`^(\S*)\s+(\S*)\s+(\S*)\s+([\S\s]+)\s*$`)
	match := regex.FindStringSubmatch(dev_line)
	if len(match) != 5 {
		fmt.Println("Error. Incorrect number of fields returned. Aborting.")
	}

	return Device{
		device:      strings.TrimSpace(match[1]),
		device_type: strings.TrimSpace(match[2]),
		state:       strings.TrimSpace(match[3]),
		conn:        strings.TrimSpace(match[4]),
	}
}
