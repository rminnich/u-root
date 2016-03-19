// Copyright 2012 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// deplist finds the names of all the files we need to include for a given
// u-root initramfs. It uses the names of all the commands.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/u-root/u-root/uroot"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	urootPath = "src/github.com/u-root/u-root"
)

type GoDirs struct {
	Dir     string
	Deps    []string
	GoFiles []string
	SFiles  []string
	Goroot  bool
}

type GoFiles struct {
	Dir        string
	ImportPath string
	Name       string
	Goroot     bool
	Standard   bool
	Root       string
	GoFiles    []string
	SFiles     []string
	Imports    []string
	Deps       []string
}

var (
	config struct {
		Debug       bool
		Goos        string
		Goroot      string
		Godotdot    string
		Arch        string
		Gopath      string
		Urootpath   string
		Fail        bool
		Dirs        map[string]bool
		Deps        map[string]bool
		GorootFiles map[string]bool
		UrootFiles  map[string]bool
	}
	debug = nodebug
)

func nodebug(string, ...interface{}) {}

func guessgoarch() {
	config.Arch = os.Getenv("GOARCH")
	if config.Arch != "" {
		config.Arch = path.Clean(config.Arch)
		return
	}
	debug("GOARCH is not set, trying to guess")
	u, err := uroot.Uname()
	if err != nil {
		debug("uname failed, using default amd64")
		config.Arch = "amd64"
	} else {
		switch {
		case u.Machine == "i686" || u.Machine == "i386" || u.Machine == "x86":
			config.Arch = "386"
		case u.Machine == "x86_64" || u.Machine == "amd64":
			config.Arch = "amd64"
		case u.Machine == "armv7l" || u.Machine == "armv6l":
			config.Arch = "arm"
		case u.Machine == "ppc" || u.Machine == "ppc64":
			config.Arch = "ppc64"
		default:
			debug("Unrecognized arch")
			config.Fail = true
		}
	}
}
func guessgoroot() {
	config.Goroot = os.Getenv("GOROOT")
	if config.Goroot != "" {
		config.Goroot = path.Clean(config.Goroot)
		debug("Using %v from the environment as the GOROOT", config.Goroot)
		config.Godotdot = path.Dir(config.Goroot)
		return
	}
	debug("Goroot is not set, trying to find a go binary")
	p := os.Getenv("PATH")
	paths := strings.Split(p, ":")
	for _, v := range paths {
		g := path.Join(v, "go")
		if _, err := os.Stat(g); err == nil {
			config.Goroot = path.Dir(path.Dir(v))
			config.Godotdot = path.Dir(config.Goroot)
			debug("Guessing that goroot is %v from $PATH", config.Goroot)
			return
		}
	}
	debug("GOROOT is not set and can't find a go binary in %v", p)
	config.Fail = true
}

func guessgopath() {
	defer func() {
		config.Godotdot = path.Dir(config.Goroot)
	}()
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		config.Gopath = gopath
		config.Urootpath = path.Join(gopath, urootPath)
		return
	}
	// It's a good chance they're running this from the u-root source directory
	log.Fatalf("Fix up guessgopath")
	cwd, err := os.Getwd()
	if err != nil {
		debug("GOPATH was not set and I can't get the wd: %v", err)
		config.Fail = true
		return
	}
	// walk up the cwd until we find a u-root entry. See if cmds/init/init.go exists.
	for c := cwd; c != "/"; c = path.Dir(c) {
		if path.Base(c) != "u-root" {
			continue
		}
		check := path.Join(c, "cmds/init/init.go")
		if _, err := os.Stat(check); err != nil {
			//debug("Could not stat %v", check)
			continue
		}
		config.Gopath = c
		debug("Guessing %v as GOPATH", c)
		os.Setenv("GOPATH", c)
		return
	}
	config.Fail = true
	debug("GOPATH was not set, and I can't see a u-root-like name in %v", cwd)
	return
}

func goListPkg(name string) (*GoDirs, error) {
	cmd := exec.Command("go", "list", "-json", name)
	if config.Debug {
		debug("Run %v @ %v", cmd, cmd.Dir)
	}
	j, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var p GoDirs
	if err := json.Unmarshal([]byte(j), &p); err != nil {
		return nil, err
	}

	debug("%v, %v %v", p, p.GoFiles, p.SFiles)
	for _, v := range append(p.GoFiles, p.SFiles...) {
		if p.Goroot {
			config.GorootFiles[path.Join(p.Dir, v)] = true
		} else {
			config.UrootFiles[path.Join(p.Dir, v)] = true
		}
	}

	return &p, nil
}

// sad news. If I concat the Go cpio with the other cpios, for reasons I don't understand,
// the kernel can't unpack it. Don't know why, don't care. Need to create one giant cpio and unpack that.
// It's not size related: if the go archive is first or in the middle it still fails.
func main() {
	config.Dirs = make(map[string]bool)
	config.Deps = make(map[string]bool)
	config.GorootFiles = make(map[string]bool)
	config.UrootFiles = make(map[string]bool)

	flag.BoolVar(&config.Debug, "d", false, "Enable debug prints")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"github.com/u-root/u-root/cmds/cat"}
	}
	if config.Debug {
		debug = log.Printf
	}

	guessgoarch()
	config.Goos = "linux"
	guessgoroot()
	guessgopath()
	if config.Fail {
		log.Fatal("Setup failed")
	}

	// It would be nice to run go list -json with lots of package names but it produces invalid JSON.
	// It produces a stream thatis {}{}{} at the top level and the decoders don't like that.
	// TODO: fix it later. Maybe use template after all. For now this is more than adequate.
	for _, v := range args {
		p, err := goListPkg(v)
		if err != nil {
			log.Fatalf("%v", err)
		}
		debug("cmd p is %v", p)
		for _, v := range p.Deps {
			config.Deps[v] = true
		}
	}

	for v := range config.Deps {
		if _, err := goListPkg(v); err != nil {
			log.Fatalf("%v", err)
		}
	}
	for v := range config.GorootFiles {
		fmt.Println(v)
	}
	for v := range config.UrootFiles {
		fmt.Println(v)
	}
}
