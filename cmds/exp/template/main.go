// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

// template walks the arg list, expanding globs as needed, and outputs a sorted list
// of matching names.
func template(names ...string) ([]string, error) {
	var visited = map[string]bool{}
	var out = map[string]struct{}{}
	for len(names) > 0 {
		var h string
		h, names = names[0], names[1:]
		if _, ok := visited[h]; ok {
			continue
		}
		visited[h] = true
		if t, ok := templates[h]; ok {
			names = append(names, t...)
		}
		l, err := filepath.Glob(h)
		if err != nil {
			return nil, err
		}
		if len(l) == 1 {
			out[l[0]] = struct{}{}
			continue
		}
		names = append(names, l...)
	}

	var sorted sort.StringSlice
	for v := range out {
		sorted = append(sorted, v)
	}
	sorted.Sort()
	return sorted, nil
}

// This dead-simple command has no options, as it is intended to be used in scripts
// or via command line invocation:
// some-command `template this that other`
func main() {
	l, err := template(os.Args[:]...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", l)
}
