package nmcli

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// Radio command type
type RadioCommand int

// Enum for Radio commands
const (
	OFF RadioCommand = iota
	ON
)

func (c RadioCommand) String() string {
	switch c {
	case 0:
		return "off"
	case 1:
		return "on"
	default:
		return ""
	}
}

// Define Radio Types
type RadioType int

const (
	WIFI RadioType = iota
	WWAN
	ALL
)

func (r RadioType) String() string {
	switch r {
	case 0:
		return "wifi"
	case 1:
		return "wwan"
	case 2:
		return "all"
	default:
		return ""
	}
}

// Shows software and hardware state of Wifi, Wwan radios.
type RadioList struct {
	Wifi   string
	WifiHW string
	Wwan   string
	WwanHW string
}

func Radios() (RadioList, error) {
	res, err := exec.Command(
		"bash",
		"-c",
		"nmcli radio",
	).Output()
	if err != nil {
		return RadioList{}, err
	}
	// Get status line
	input := strings.Split(strings.TrimSpace(string(res[:])), "\n")[1]
	// parse
	return parseRadios(input), nil
}

func (r *RadioList) ChangeState(radio RadioType, state RadioCommand) (string, error) {
	res, err := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf("nmcli radio %s %s", radio, state),
	).Output()
	// update radio list
	*r, _ = Radios()
	return string(res), err
}

//*********************
// HELPERS
// ********************

// Parses the status line of nmcli radio.
// Command is always returned in this order: WifiHW, Wifi, WwanHW, Wwan
func parseRadios(radio_states string) RadioList {
	regex := regexp.MustCompile(`^(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s*`)
	match := regex.FindStringSubmatch(radio_states)
	if len(match) != 5 {
		fmt.Println("Error. Incorrect number of fields returned. Aborting.")
		return RadioList{}
	}
	return RadioList{
		WifiHW: strings.TrimSpace(match[1]),
		Wifi:   strings.TrimSpace(match[2]),
		WwanHW: strings.TrimSpace(match[3]),
		Wwan:   strings.TrimSpace(match[4]),
	}
}
