// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !linux

package pci

import "os"

// OnePCI takes the name of a directory containing linux-style
// PCI files and returns a filled-in *PCI.
func OnePCI(dir string) (*PCI, error) {
	return nil, os.ErrNotExist
}

// BaseLimType parses a Linux resource string into base, limit, and attributes.
// The string must have three hex fields.
// Gaul was divided into three parts.
// So are the BARs.
func BaseLimType(_ string) (uint64, uint64, uint64, error) {
	return 0, 0, 0, os.ErrInvalid
}
