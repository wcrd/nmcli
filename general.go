package nmcli

import (
	"errors"
	"os/exec"
	"regexp"
	"strings"
)

type State struct {
	state        string
	connectivity string
	wifi_hw      bool
	wifi         bool
	wwan_hw      bool
	wwan         bool
}

func Status() (state State, err error) {

	res, err := exec.Command(
		"bash",
		"-c",
		"nmcli general status",
	).CombinedOutput()

	if err != nil {
		return State{}, err
	}

	// parse result into useful struct
	// remove header
	statuses := strings.Split(string(res), "\n")[1]
	// extract individual statuses
	regex := regexp.MustCompile(`^([\S\s]+)\s{2}(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s*`)
	match := regex.FindStringSubmatch(statuses)
	if len(match) != 7 {
		return State{}, errors.New("incorrect number of fields returned...aborting")
	}

	// process into struct
	return State{
		state:        match[1],
		connectivity: match[2],
		wifi_hw:      match[3] == "enabled",
		wifi:         match[4] == "enabled",
		wwan_hw:      match[5] == "enabled",
		wwan:         match[6] == "enabled",
	}, nil

}
