// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// testramfs tests things, badly
package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"syscall"

	flag "github.com/spf13/pflag"

	"github.com/u-root/u-root/pkg/cpio"
	"github.com/u-root/u-root/pkg/pty"
)

const (
	unshareFlags = syscall.CLONE_NEWNS
	cloneFlags   = syscall.CLONE_NEWIPC |
		syscall.CLONE_NEWNET |
		// making newpid work will be more tricky,
		// since none of my CLs to fix go runtime for
		// it ever got in.
		//syscall.CLONE_NEWPID |
		syscall.CLONE_NEWUTS |
		0
)

var (
	noremove    = flag.BoolP("noremove", "n", false, "remove tempdir when done")
	interactive = flag.BoolP("interactive", "i", false, "interactive mode")
)

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatalln("usage: %s <cpio-path>", os.Args[0])
	}

	c := flag.Args()[0]

	f, err := os.Open(c)
	if err != nil {
		log.Fatal(err)
	}

	// So, what's the plan here?
	//
	// - new mount namespace
	//   - root mount is a tmpfs mount filled with the archive.
	//
	// - new PID namespace
	//   - archive/init actually runs as PID 1.

	tempDir, err := ioutil.TempDir("", "u-root")
	if err != nil {
		log.Fatal(err)
	}
	// Don't do a RemoveAll. This should be empty and
	// an error can tell us we got something wrong.
	if !*noremove {
		defer func(n string) {
			log.Printf("Removing %v", n)
			if err := os.Remove(n); err != nil {
				log.Fatal(err)
			}
		}(tempDir)
	}
	if err := syscall.Mount("", tempDir, "tmpfs", 0, ""); err != nil {
		log.Fatal(err)
	}
	if !*noremove {
		defer func(n string) {
			log.Printf("Unmounting %v", n)
			if err := syscall.Unmount(n, syscall.MNT_DETACH); err != nil {
				log.Fatal(err)
			}
		}(tempDir)
	}

	archiver, err := cpio.Format("newc")
	if err != nil {
		log.Fatal(err)
	}

	r := archiver.Reader(f)
	for {
		rec, err := r.ReadRecord()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		cpio.CreateFileInRoot(rec, tempDir)
	}

	cmd, err := pty.New()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Command("/init")
	cmd.C.SysProcAttr.Chroot = tempDir
	cmd.C.SysProcAttr.Cloneflags = cloneFlags
	cmd.C.SysProcAttr.Unshareflags = cloneFlags
	if *interactive {
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
		return
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	go io.Copy(cmd.TTY, cmd.Ptm)

	// At this point you could use an array of commands/output templates to
	// drive the test, and end with the exit command shown nere.
	if n, err := cmd.Ptm.Write([]byte("exit\n")); err != nil {
		log.Printf("Writing exit: want (5, nil); got (%d, %v)\n", n, err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

}
