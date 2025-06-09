// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// extract a tar file, to a tempdir, and run a command in it with args.
//
// Synopsis:
//
//	runcontainer [OPTION...] image command [args]
//
// Description:
//
//	This command extracts a tar file to a directory, and runs
//	a command with args chroot'ed in that directory.
//
// Options:
//
//	-v: verbose
package main

import (
	"archive/tar"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/u-root/u-root/pkg/tarutil"
	"github.com/u-root/u-root/pkg/uroot/unixflag"
)

var ErrUsage = errors.New("usage")

type cmd struct {
	p     params
	where string
	image string
	argv0 string
	args  []string
}

type params struct {
	verbose bool
}

// SafeFilter filters out all files which are not regular, symlinks,  and not directories.
// It also sets appropriate permissions.
func SafeFilter(hdr *tar.Header) bool {
	if hdr.Typeflag == tar.TypeDir {
		hdr.Mode = 0o770
		return true
	}
	if hdr.Typeflag == tar.TypeReg {
		hdr.Mode = 0o660
		return true
	}
	if hdr.Typeflag == tar.TypeLink {
		hdr.Mode = 0o660
		return true
	}
	if hdr.Typeflag == tar.TypeSymlink {
		hdr.Mode = 0o777
		return true
	}
	return false
}
func command(p params, args []string) (*cmd, error) {
	if len(args) < 2 {
		return nil, ErrUsage
	}
	d, err := ioutil.TempDir("", "runcontainer")
	if err != nil {
		return nil, err
	}
	return &cmd{
		p:     p,
		where: d,
		image: args[0],
		argv0: args[1],
		args:  args[2:],
	}, nil
}

func (c *cmd) run() error {
	opts := &tarutil.Opts{Filters: []tarutil.Filter{SafeFilter}}
	if c.p.verbose {
		opts.Filters = append(opts.Filters, tarutil.VerboseFilter)
	}

	f, err := os.Open(c.image)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := tarutil.ExtractDir(f, c.image, opts); err != nil {
		return err
	}
	return nil
}

func main() {
	var (
		verbose bool
	)
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	f.BoolVar(&verbose, "v", false, "print each filename")

	f.Parse(unixflag.OSArgsToGoArgs())
	cmd, err := command(params{verbose: verbose}, f.Args())
	if err != nil {
		f.Usage()
		log.Fatal(err)
	}
	if err := cmd.run(); err != nil {
		log.Fatal(err)
	}
}
