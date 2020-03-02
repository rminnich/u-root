// Copyright 2013-2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// madeye merges u-root initramfs to form a single universal initramfs.
// Synopsis:
//     madeye initramfs [initramfs...]
//
// Description:
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/uio"
	"golang.org/x/sys/unix"
)

var (
	debug = func(string, ...interface{}) {}
	d     = flag.Bool("v", false, "Debug prints")
	arch  = map[string]string{
		"initramfs.linux_amd64.cpio": "amd64",
		"initramfs.linux_arm.cpio":   "arm",
	}
	out = map[string]*cpio.Record{}
)

func usage() {
	log.Fatalf("Usage: madeye")
}

func file(archiver cpio.RecordFormat, n string, f io.ReaderAt) ([]cpio.Record, error) {
	var r []cpio.Record
	rr := archiver.Reader(f)
	a, ok := arch[filepath.Base(n)]
	if !ok {
		return r, fmt.Errorf("%s: don't know about this", n)
	}
	debug("arch is %s", a)
	for {
		rec, err := rr.ReadRecord()
		if err == io.EOF {
			break
		}
		debug("Read %v", rec)
		if err != nil {
			log.Fatalf("error reading records: %v", err)
		}
		d := filepath.Dir(rec.Name)
		switch d {
		case "bbin", "bin":
			rec.Name = filepath.Join(a, rec.Name)
			debug("Change to %v", rec)
		default:
			debug("dir is %v, ignore", d)
		}
		switch rec.Name {
		case "init", "bbin", "bin":
			rec.Name = filepath.Join(a, rec.Name)
		}
		// TODO: make this use os constants, not unix constants.
		switch rec.Mode & unix.S_IFMT {
		case unix.S_IFLNK:
			content, err := ioutil.ReadAll(uio.Reader(rec))
			if err != nil {
				return nil, err
			}
			switch string(content) {
			case "bbin", "bin":
				content = []byte(filepath.Join(a, string(content)))
				debug("Change to %v", rec)
			default:
				debug("dir is %v, ignore", d)
			}
			rec.ReaderAt = bytes.NewReader(content)
		}

		if _, ok := out[rec.Name]; ok {
			continue
		}
		out[rec.Name] = &rec
		r = append(r, rec)
	}
	return r, nil
}
func main() {
	flag.Parse()
	if *d {
		debug = log.Printf
	}

	a := flag.Args()
	debug("Args %v", a)
	if len(a) < 1 {
		usage()
	}

	archiver, err := cpio.Format("newc")
	if err != nil {
		log.Fatal(err)
	}

	var rr []cpio.Record
	for _, a := range flag.Args() {
		f, err := os.Open(a)
		if err != nil {
			log.Fatal(err)
		}
		r, err := file(archiver, a, f)
		if err != nil {
			log.Fatal(err)
		}
		// Why not a defer? Because that would happen
		// outside the for loop. Not that it really matters:
		// any kind of explicit close is a bit silly here, we're
		// never going to have more than MAXFD arguments anyway,
		// but better safe than sorry.
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
		rr = append(rr, r...)
	}
	// process ...
	archiver, err = cpio.Format("newc")
	if err != nil {
		log.Fatal(err)
	}
	rw := archiver.Writer(os.Stdout)
	for _, r := range rr {
		if *d {
			log.Printf("%s", r)
			continue
		}
		if err := rw.WriteRecord(r); err != nil {
			log.Fatal(err)
		}
	}
	if !*d {
		if err := cpio.WriteTrailer(rw); err != nil {
			log.Fatalf("Error writing trailer record: %v", err)
		}
	}
}
