// Copyright 2009-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ldd returns all the library dependencies of an executable.
// See note below on why we do not really do this on OSX.
package ldd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Follow starts at a pathname and adds it
// to a map if it is not there.
// If the pathname is a symlink, indicated by the Readlink
// succeeding, links repeats and continues
// for as long as the name is not found in the map.
func follow(l string, names map[string]*FileInfo) error {
	for {
		if names[l] != nil {
			return nil
		}
		i, err := os.Lstat(l)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

		names[l] = &FileInfo{FullName: l, FileInfo: i}
		if i.Mode().IsRegular() {
			return nil
		}
		// If it's a symlink, the read works; if not, it fails.
		// we can skip testing the type, since we still have to
		// handle any error if it's a link.
		next, err := os.Readlink(l)
		if err != nil {
			return err
		}
		// It may be a relative link, so we need to
		// make it abs.
		if filepath.IsAbs(next) {
			l = next
			continue
		}
		l = filepath.Join(filepath.Dir(l), next)
	}
}

// dyn runs the binary with DYLD_PRINT_LIBRARIES=1.
// This is insanely dangerous.
// What if the binary is ... well, anyway, this is insanely dangerous.
// OSX is stupid.
// So we're going to return the name, and that's it, until we can find
// a way to do this without running the command.
func dyn(n string) ([]string, error) {
	return []string{n}, nil
}

type FileInfo struct {
	FullName string
	os.FileInfo
}

// Ldd returns a list of all library dependencies for a set of files.
//
// If a file has no dependencies, that is not an error. The only possible error
// is if a file does not exist, or it says it has an interpreter but we can't
// read it, or we are not able to run its interpreter.
//
// It's not an error for a file to not be an ELF.
func Ldd(names []string) ([]*FileInfo, error) {
	var (
		list    = make(map[string]*FileInfo)
		libs    []*FileInfo
	)
	for _, n := range names {
		if err := follow(n, list); err != nil {
			return nil, err
		}
	}
	for _, n := range names {
		n, err := dyn(n)
		if err != nil {
			return nil, err
		}
		for i := range n {
			if err := follow(n[i], list); err != nil {
				log.Fatalf("ldd: %v", err)
			}
		}
	}

	for i := range list {
		libs = append(libs, list[i])
	}

	return libs, nil
}

// List returns the dependency file paths of files in names.
func List(names []string) ([]string, error) {
	var list []string
	l, err := Ldd(names)
	if err != nil {
		return nil, err
	}
	for i := range l {
		list = append(list, l[i].FullName)
	}
	return list, nil
}
