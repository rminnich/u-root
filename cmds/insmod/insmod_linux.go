// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Insert a module into the Linux kernel
//
// Synopsis:
//	insmod [filename] [module options...]
//
// Description:
//	insmod is a clone of insmod(8)
package main

import (
	"log"
	"os"
	"strings"

	"github.com/u-root/u-root/pkg/kmodule"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("insmod: ERROR: missing filename.\n")
	}

	// get filename from argv[1]
	filename := os.Args[1]

	// Everything else is module options
	options := strings.Join(os.Args[2:], " ")

	if err := kmodule.LoadModule(filename, options, 0, kmodule.ProbeOpts{}); err != nil {
		log.Fatalf("insmod: could not load %q: %v", filename, err)
	}
}
