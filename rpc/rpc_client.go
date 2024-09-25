package rpc

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"log"

	"github.com/lwlcom/cisco_exporter/connector"
)

const (
	IOSXE string = "IOSXE"
	NXOS  string = "NXOS"
	IOS   string = "IOS"
)

// Client sends commands to a Cisco device
type Client struct {
	conn       *connector.SSHConnection
	Debug      bool
	OSType     string
	interfaces []string
}

// NewClient creates a new client connection
func NewClient(ssh *connector.SSHConnection, debug bool) *Client {
	rpc := &Client{conn: ssh, Debug: debug}

	return rpc
}

// Identify tries to identify the OS running on a Cisco device
func (c *Client) Identify() error {
	output, err := c.RunCommand("show version")
	if err != nil {
		return err
	}
	switch {
	case strings.Contains(output, "IOS XE"):
		c.OSType = IOSXE
	case strings.Contains(output, "IOS-XE"):
		c.OSType = IOSXE
	case strings.Contains(output, "NX-OS"):
		c.OSType = NXOS
	case strings.Contains(output, "IOS Software"):
		c.OSType = IOS
	default:
		return errors.New("Unknown OS")
	}
	if c.Debug {
		log.Printf("Host %s identified as: %s\n", c.conn.Host, c.OSType)
	}
	return nil
}

// RunCommand runs a command on a Cisco device
func (c *Client) RunCommand(cmd string) (string, error) {
	if c.Debug {
		log.Printf("Running command on %s: %s\n", c.conn.Host, cmd)
	}
	output, err := c.conn.RunCommand(fmt.Sprintf("%s", cmd))
	if err != nil {
		println(err.Error())
		return "", err
	}

	return output, nil
}

// Runs command to show interfaces and returns list of interface names
func (c *Client) GetInterfaceNames(includeVirtual bool) ([]string, error) {
	if c.OSType != IOSXE && c.OSType != NXOS && c.OSType != IOS {
		return nil, errors.New("'show interfaces stats' is not implemented for " + c.OSType)
	}
	var items []string
	if len(c.interfaces) < 1 {
		err := c.PopulateInterfaces()
		if err != nil {
			return items, err
		}
	}
	virtualNames := [4]string{"Vlan", "Loopback", "Tunnel", "Port-channel"}
IFACE:
	for _, name := range c.interfaces {
		if !includeVirtual {
			// ignore virtual interfaces
			for _, virtualName := range virtualNames {
				if strings.HasPrefix(name, virtualName) {
					continue IFACE
				}
			}
		}
		items = append(items, name)
	}
	return items, nil
}

func (c *Client) PopulateInterfaces() error {
	var items []string
	out, err := c.RunCommand("show interfaces stats")
	if err != nil {
		return err
	}
	deviceNameRegexp, _ := regexp.Compile(`^(?:Interface\s)?([a-zA-Z0-9\/\.-]+)(?: is disabled)?\s*$`)
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		matches := deviceNameRegexp.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		items = append(items, matches[1])
	}
	c.interfaces = items
	return nil
}
