package nmcli

import (
	"errors"
	"fmt"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
)

type Connection struct {
	Name      string
	Uuid      string
	Conn_type string
	Device    string // aka ifname
	Addr      *AddressDetails
}

// Deletes the connection.
// Returns nmcli success message and error.
func (c Connection) Delete() (string, error) {
	res, err := exec.Command("nmcli", "connection", "del", c.Name).Output()
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// Modifies the connection with given parameters.
// func (c Connection) Modify() (string, error)

type NewConnectionDetails struct {
	Name      string
	Conn_type string
	Ifname    string
	Addr      *AddressDetails
}

type AddressDetails struct {
	Ipv4_method  string   `cmd:"ipv4.method"`
	Ipv4_address string   `cmd:"ipv4.address"`
	Ipv4_gateway string   `cmd:"ipv4.gateway"`
	Ipv4_dns     []string `cmd:"ipv4.dns"`
}

// Returns all connections defined in nmcli
// Equivalent to: nmcli connection
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

// Finds connection by con-name, if it exists.
// Returns a list of Connections and an error.
// Equivalent to: nmcli connection show {name}
func GetConnectionByName(conn string) ([]Connection, error) {
	// get connections
	conns := Connections()
	// check if connection with name exists
	existingConns := []Connection{}
	for _, c := range conns {
		if c.Name == conn {
			existingConns = append(existingConns, c)
		}
	}
	// single conn = OK, multi conn = ERROR, no conn = ERROR
	switch len(existingConns) {
	case 0:
		return existingConns, errors.New("no connection found")
	case 1:
		return existingConns, nil
	default:
		return existingConns, errors.New("multiple connections found")
	}
}

// Creates a new connection
// Equivalent to: nmcli con add con-name {name} type {type} ifname {ifname}
// Returns nmcli message and error
func AddConnection(conn *Connection) (string, error) {
	// Create new connection
	// TODO: Is it worth doing this in two parts? Or should execute as one command?
	res, err := exec.Command("nmcli", "connection", "add", "con-name", conn.Name, "type", conn.Conn_type, "ifname", conn.Device).Output()
	if err != nil {
		return string(res), err
	}

	// Update connection with address details
	cmds := conn.Addr.construct_commands()
	fmt.Println(append([]string{"connection", "mod", conn.Name}, cmds...))
	res, err = exec.Command("nmcli", append([]string{"connection", "mod", conn.Name}, cmds...)...).CombinedOutput()
	if err != nil {
		return string(res), err
	}
	return string(res), nil
}

func (addr *AddressDetails) construct_commands() []string {
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
		Name:      strings.TrimSpace(match[1]),
		Uuid:      strings.TrimSpace(match[2]),
		Conn_type: strings.TrimSpace(match[3]),
		Device:    strings.TrimSpace(match[4]),
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
