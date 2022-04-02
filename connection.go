package nmcli

import (
	"fmt"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
)

type Connection struct {
	name      string
	uuid      string
	conn_type string
	device    string
}

type NewConnectionDetails struct {
	Name      string
	Conn_type string
	Ifname    string
	Addr      *NewAddressDetails
}

type NewAddressDetails struct {
	Ipv4_method  string   `cmd:"ipv4.method"`
	Ipv4_address string   `cmd:"ipv4.address"`
	Ipv4_gateway string   `cmd:"ipv4.gateway"`
	Ipv4_dns     []string `cmd:"ipv4.dns"`
}

func Connections() []Connection {
	res, err := exec.Command("nmcli", "connection").Output()
	if err != nil {
		// handle error
		fmt.Println(err)
	}
	// process result
	results := make([]Connection, 0)
	input := strings.Split(strings.TrimSpace(string(res[:])), "\n")
	// fmt.Printf("%+v\n", input)
	// pop first row (headers)
	for _, line := range input[1:] {
		// fmt.Println(line)
		results = append(results, parseConnection(line))
	}

	return results
}

func AddConnection(conn *NewConnectionDetails) error {
	// validate input
	//		check ip vs subnet, if ip/24 etc.

	// Create new connection
	// TODO: Is it worth doing this in two parts? Or should execute as one command?
	_, err := exec.Command("nmcli", "connection", "add", "con-name", conn.Name, "type", conn.Conn_type, "ifname", conn.Ifname).Output()
	if err != nil {
		// handle error
		fmt.Println(err)
		return err
	}

	// Update connection with address details
	cmds := conn.Addr.construct_commands()
	fmt.Println(append([]string{"connection", "mod", conn.Name}, cmds...))
	res, err := exec.Command("nmcli", append([]string{"connection", "mod", conn.Name}, cmds...)...).CombinedOutput()
	if err != nil {
		// handle error
		fmt.Println(err)
		fmt.Printf("%s\n", res)
		return err
	}
	return nil
}

func (addr *NewAddressDetails) construct_commands() []string {
	output := make([]string, 0)
	// Get type
	t := reflect.TypeOf(*addr)

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Field(i)

		// Get the field tag value
		tag := field.Tag.Get("cmd")

		// Get the field value
		value := reflect.ValueOf(*addr).Field(i)

		// if not empty, write command
		if !value.IsZero() {
			switch x := value.Interface().(type) {
			case string:
				output = append(output, []string{tag, value.String()}...)
			case []string:
				output = append(output, []string{tag, fmt.Sprintf("%v", strings.Join(value.Interface().([]string), " "))}...)
			default:
				fmt.Println(x)
			}
		}

	}
	return output
}

// nmcli connection add con-name {name} type {type} ifname {device}

//*********************
// HELPERS
// ********************

func parseConnection(conn_line string) Connection {
	regex := regexp.MustCompile(`^([\S\s]+)\s{2}(\S+)\s{2}(\S+)\s+(\S+)\s*`)
	match := regex.FindStringSubmatch(conn_line)
	// fmt.Println(match)
	if len(match) != 5 {
		fmt.Println("Error. Incorrect number of fields returned. Aborting.")
	}

	return Connection{
		name:      strings.TrimSpace(match[1]),
		uuid:      strings.TrimSpace(match[2]),
		conn_type: strings.TrimSpace(match[3]),
		device:    strings.TrimSpace(match[4]),
	}
}

func (conn *NewConnectionDetails) fill_defaults() {
	// Fills default values in for new connection creation
	if conn.Name == "" {
		conn.Name = "new-con-" + conn.Ifname
	}
}

func Test() {
	res, _ := exec.Command("nmcli", []string{"device", "show"}...).Output()
	fmt.Println(res)
}
