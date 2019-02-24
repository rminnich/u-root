// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// crosboot loads a chromeos kernel partition image,
// pulls it apart, then kexecs it.
//
// Synopsis:
//     crosboot <file or partition>
//
// Description:
//     Loads and executes a ChromeOS kernel from a ChromeOS partition or image.
//
// Options:
package main

import (
	"io/ioutil"
	"log"

	flag "github.com/spf13/pflag"
	"github.com/u-root/u-root/pkg/kexec"
)

type croskernel struct {
	preamble []byte
	kern     []byte
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatalf("Usage: crosboot <file>")
	}
	f := flag.Arg(0)

	b, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}

	k := &croskernel{preamble: b[:64*1024], kern: b[64*1024:]}
	s := kexec.NewSegment(k.kern[64*1024:], kexec.Range{
		Start: uintptr(0x100000),
		Size:  uint(len(k.kern[64*1024:])),
	})
	m := &kexec.Memory{}
	if err := m.ParseMemoryMap(); err != nil {
		log.Fatal(err)
	}
	m.Segments = append(m.Segments, s)

	// This is stupid but for
	/*
		if err := kexec.Load(m.EntryPoint, m.Segments(), 0); err != nil {
			log.Fatalf("kexec.Load() error: %v", err)
		}
		if err := kexec.Reboot(); err != nil {
			log.Fatalf("kexec.Reboot() error: %v", err)
		}
	*/
}
