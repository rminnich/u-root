// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"bytes"
	"fmt"
	"net"
)

// gencsum generates a uint8 checksum of a []uint8
func gencsum(b []uint8) uint8 {
	var csum uint8
	for _, bb := range b {
		csum += bb
	}
	Debug("csum %#x %#x across %d bytes", csum, ^csum, len(b))
	return ^csum + 1
}

// MarshalBasicTypes marshals simple non-struct types into the head and heap.
func MarshalBasicTypes(head, heap *bytes.Buffer, i interface{}) error {
	switch s := i.(type) {
	case sockaddr:
		Debug("addr")
		a, err := net.ResolveTCPAddr("tcp", string(s))
		if err != nil {
			return fmt.Errorf("addr %s: %v", s, err)
		}
		w(head, a.IP.To16(), uint16(a.Port))
	case ipaddr:
		a, err := net.ResolveIPAddr("ip", string(s))
		if err != nil {
			return fmt.Errorf("addr %s: %v", s, err)
		}
		w(head, a.IP.To16())
		Debug("net")
	case flag:
	case mac:
		hw, err := net.ParseMAC(string(s))
		if err != nil {
			return err
		}
		if len(hw) != 6 {
			return fmt.Errorf("%q is not an ethernet MAC", s)
		}
		w(head, hw)
		Debug("mac")
	case bdf:
		if err := uw(head, string(s), 16); err != nil {
			return err
		}
		Debug("bdf")
	case u8:
		if err := uw(head, string(s), 8); err != nil {
			return err
		}

	case u16:
		if err := uw(head, string(s), 16); err != nil {
			return err
		}

	case u64:
		if err := uw(head, string(s), 64); err != nil {
			return err
		}
	case sheap:
		w(head, uint16(len(s)), uint16(heap.Len()))
		Debug("Write %q to heap", string(s))
		w(heap, []byte(s))
	default:
		return fmt.Errorf("Don't know what to do with %T", s)
	}
	return nil
}
