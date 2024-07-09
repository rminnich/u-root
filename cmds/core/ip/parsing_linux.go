// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/vishvananda/netlink"
)

func findPrefix(cmd string, cmds []string) string {
	var x, n int

	for i, v := range cmds {
		if strings.HasPrefix(v, cmd) {
			n++
			x = i
		}
	}

	if n == 1 {
		return cmds[x]
	}

	return ""
}

var ErrNotFound = fmt.Errorf("not found")

// in the ip command, turns out 'dev' is a noise word.
// The BNF it shows is not right in that case.
// Always make 'dev' optional.
func parseDeviceName(mandatory bool) (netlink.Link, error) {
	switch mandatory {
	case true:
		cursor++
		whatIWant = []string{"dev", "device name"}

		if arg[cursor] == "dev" {
			cursor++
		}

		whatIWant = []string{"device name"}
		return netlink.LinkByName(arg[cursor])
	case false:
		if cursor == len(arg)-1 {
			return nil, ErrNotFound
		}

		cursor++
		whatIWant = []string{"dev", "device name"}

		if cursor > len(arg)-1 {
			return nil, ErrNotFound
		}

		if arg[cursor] == "dev" {
			cursor++

			if cursor > len(arg)-1 {
				return nil, ErrNotFound
			}

		}

		whatIWant = []string{"device name"}
		return netlink.LinkByName(arg[cursor])
	}

	return nil, ErrNotFound
}

func parseType() (string, error) {
	if cursor == len(arg)-1 {
		return "", ErrNotFound
	}

	cursor++
	whatIWant = []string{"type"}

	if cursor > len(arg)-1 {
		return "", ErrNotFound
	}

	if arg[cursor] != "type" {
		return "", ErrNotFound
	}

	cursor++

	whatIWant = []string{"type name"}
	return arg[cursor], nil
}

func parseHardwareAddress() (net.HardwareAddr, error) {
	cursor++
	whatIWant = []string{"<mac address>"}

	return net.ParseMAC(arg[cursor])
}

func parseString() string {
	cursor++
	whatIWant = []string{"string"}

	return arg[cursor]
}

func parseInt() (int, error) {
	cursor++
	whatIWant = []string{"<id>"}

	return strconv.Atoi(arg[cursor])
}

func parseUint32() (uint32, error) {
	cursor++
	whatIWant = []string{"<id>"}

	val, err := strconv.ParseUint(arg[cursor], 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse uint32: %v", err)
	}

	return uint32(val), nil
}

func parseBool() (bool, error) {
	cursor++
	whatIWant = []string{"true", "false"}

	switch arg[cursor] {
	case "on":
		return true, nil
	case "off":
		return false, nil
	}

	return false, fmt.Errorf("invalid bool value: %v", arg[cursor])
}

func parseName() (string, error) {
	cursor++
	whatIWant = []string{"name", "device name"}
	if arg[cursor] == "name" {
		cursor++
	}

	whatIWant = []string{"device name"}

	return arg[cursor], nil
}

func parseNodeSpec() string {
	cursor++
	whatIWant = []string{"default", "CIDR"}

	return arg[cursor]
}

func parseNextHop() (string, net.IP, error) {
	cursor++
	whatIWant = []string{"via"}

	if arg[cursor] != "via" {
		return "", nil, usage()
	}

	nh := arg[cursor]
	cursor++
	whatIWant = []string{"Gateway CIDR"}

	addr := net.ParseIP(arg[cursor])
	if addr == nil {
		return "", nil, fmt.Errorf("failed to parse gateway IP: %v", arg[cursor])
	}

	return nh, addr, nil
}