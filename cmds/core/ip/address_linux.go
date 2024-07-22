// Copyright 2024 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"

	"github.com/vishvananda/netlink"
)

const addressHelp = `Usage: ip address {add|replace} ADDR dev IFNAME 

       ip address del IFADDR dev IFNAME 

	   ip address flush

       ip address [ show [ dev IFNAME ] [ type TYPE ]

	   ip address help

TYPE := { bareudp | bond | bond_slave | bridge | bridge_slave |
          dummy | erspan | geneve | gre | gretap | ifb |
          ip6erspan | ip6gre | ip6gretap | ip6tnl |
          ipip | ipoib | ipvlan | ipvtap |
          macsec | macvlan | macvtap |
          netdevsim | nlmon | rmnet | sit | team | team_slave |
          vcan | veth | vlan | vrf | vti | vxcan | vxlan | wwan |
          xfrm }
`

func (cmd cmd) address() error {
	if len(arg) == 1 {
		return cmd.showAllLinks(true)
	}
	cursor++
	expectedValues = []string{"add", "replace", "del", "show", "flush", "help"}
	argument := arg[cursor]

	c := findPrefix(argument, expectedValues)
	switch c {
	case "show":
		return cmd.addressShow()
	case "add", "change", "replace", "del":
		return cmd.addressAddReplaceDel(c)
	case "flush":
		return cmd.addressFlush()
	case "help":
		fmt.Fprint(cmd.out, addressHelp)
		return nil
	default:
		return usage()
	}
}

func (cmd cmd) addressShow() error {
	device, err := parseDeviceName(false)
	if errors.Is(err, ErrNotFound) {
		return cmd.showAllLinks(true)
	}
	typeName, err := parseType()
	if errors.Is(err, ErrNotFound) {
		return cmd.showLink(device, true)
	}

	return cmd.showLink(device, true, typeName)
}

func (cmd cmd) addressAddReplaceDel(argument string) error {
	cursor++
	expectedValues = []string{"CIDR format address"}
	addr, err := netlink.ParseAddr(arg[cursor])
	if err != nil {
		return err
	}

	iface, err := parseDeviceName(true)
	if err != nil {
		return err
	}

	c := findPrefix(argument, expectedValues)
	switch c {
	case "add":
		if err := cmd.handle.AddrAdd(iface, addr); err != nil {
			return fmt.Errorf("adding %v to %v failed: %v", arg[1], arg[2], err)
		}
	case "replace":
		if err := cmd.handle.AddrReplace(iface, addr); err != nil {
			return fmt.Errorf("replacing %v on %v failed: %v", arg[1], arg[2], err)
		}
	case "del":
		if err := cmd.handle.AddrDel(iface, addr); err != nil {
			return fmt.Errorf("deleting %v from %v failed: %v", arg[1], arg[2], err)
		}
	default:
		return fmt.Errorf("subcommand %s not yet implemented, expected: %v", c, expectedValues)
	}
	return nil
}

func (cmd cmd) addressFlush() error {
	iface, err := parseDeviceName(true)
	if err != nil {
		return err
	}
	addr, err := cmd.handle.AddrList(iface, netlink.FAMILY_ALL)
	if err != nil {
		return err
	}

	for _, a := range addr {
		for idx := 1; idx <= f.loops; idx++ {
			if err := cmd.handle.AddrDel(iface, &a); err != nil {
				if idx != f.loops {
					continue
				}

				return fmt.Errorf("deleting %v from %v failed: %v", a, iface, err)
			}

			break
		}
	}

	return nil
}
