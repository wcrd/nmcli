# nmcli #

The nmcli package is a simple golang wrapper for the linux network-manager cli client (`nmcli`).

On Linux, `nmcli` is a command-line tool for controlling NetworkManager and reporting network status. It is used to create, display, edit, delete, activate, and deactivate network connections, as well as control and display network device status.

Inspired by the wonderful golang package [netlink](https://github.com/vishvananda/netlink) - the nmcli wrapper has been written to help users complete simple networking tasks in Linux that do not require the full power of the netlink package (or that may not yet be implemented).

Also inspired by the python nmcli wrapper https://github.com/ushiboy/nmcli [(pypi)](https://pypi.org/project/nmcli/).

## Dependencies ##

* NetworkManager
  * `sudo apt-get install network-manager` (Ubuntu|Debian)

## Usage ##

### Check for network-manager ###
### Radios ###
Get current state of radios:
```golang
radios, err := nmcli.Radios()
```

Change radio state:
```golang
// First get radios object
radios, err := nmcli.Radios()

// Set WIFI state
msg, err := radios.ChangeState(nmcli.WIFI, nmcli.OFF)

// Set state of all radios
msg, err := radios.ChangeState(nmcli.ALL, nmcli.ON)
```

### Connections ###

Get all connections:
```golang
list_of_connections := nmcli.Connections()
```

Get connection(s) by connection name:
```golang
list_of_connections, err := nmcli.GetConnectionByName("{con-name}")
```

Add new connection:
```golang
```

#### Connection Object ####
Delete connection:
```golang
c_list, _ := nmcli.GetConnectionByName("{con-name}")
// take first instance of connection with name {con-name}
// typically only one connection object is returned, but it is possible to have multiple connections with the same con-name
c := c_list[0]
// This will delete all connections that have name {con-name}
msg, err := c.Delete()
```

Modify connection:
```golang
c_list, _ := nmcli.GetConnectionByName("{con-name}")
// take first instance of connection with name {con-name}
// typically only one connection object is returned, but it is possible to have multiple connections with the same con-name
c := c_list[0]

// updates desired to connection details
// See docs for supported fields
c_updates := Connection{
    Name: "new-name",
    Device: "wlp58s0",
    Addr: &nmcli.AddressDetails{
      Ipv4_method:  "manual",
      Ipv4_address: "192.168.2.1",
      Ipv4_dns:     []string{"8.8.8.8", "1.1.1.1"},
    }
}
msg, err := c.Modify(c_updates)
```


### Devices ###


## Compatibility Table ##

| Object | Command | Status |
|--------|---------|--------|
| general | | not supported |
| general | status | not supported |
| general | hostname | not supported |
| general | permissions | not supported |
| general | logging | not supported |
| networking | | not supported |
| networking | on | not supported |
| networking | off | not supported |
| networking | connectivity | not supported |
| radio | | supported |
| radio | all | supported |
| radio | wifi | supported |
| radio | wwan | supported |
| connection | | supported |
| connection | show | supported |
| connection | up | not supported |
| connection | down | not supported |
| connection | add | supported |
| connection | modify | supported |
| connection | clone | not supported |
| connection | edit | not supported |
| connection | delete | supported |
| connection | reload | not supported |
| connection | load | not supported |
| connection | import | not supported |
| connection | export | not supported |
| device | | not supported |
| device | status | not supported |
| device | show | not supported |
| device | set | not supported |
| device | connect | not supported |
| device | reapply | not supported |
| device | modify | not supported |
| device | disconnect | not supported |
| device | delete | not supported |
| device | monitor | not supported |
| device | wifi | not supported |
| device | wifi connect | not supported |
| device | wifi rescan | not supported |
| device | wifi hotspot | not supported |
| device | lldp | not supported |
| agent | | not supported |
| agent | secret | not supported |
| agent | polkit | not supported |
| agent | all | not supported |
| monitor | | not supported |