// Copyright 2009-2017 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// ldd returns all the library dependencies
// of a list of file names.
// The way this is done on GNU-based systems
// is interesting. For each ELF, one finds the
// .interp section. If there is no interpreter
// there's not much to do. If there is an interpreter,
// we run it with the --list option and the file as an argument.
// We need to parse the output.
// For all lines with =>  as the 2nd field, we take the
// 3rd field as a dependency. The field may be a symlink.
// Rather than stat the link and do other such fooling around,
// we can do a readlink on it; if it fails, we just need to add
// that file name; if it succeeds, we need to add that file name
// and repeat with the next link in the chain. We can let the
// kernel do the work of figuring what to do if and when we hit EMLINK.
package uroot

import (
	"debug/elf"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Follow starts at a pathname and adds it
// to a map if it is not there.
// If the pathname is a symlink, indicated by the Readlink
// succeeding, links repeats and continues
// for as long as the name is not found in the map.
// It is not an error for a file not to exist; sometimes
// links point off into the air intentionally.
func follow(l string, names map[string]bool) {
	for {
		if names[l] {
			return
		}
		names[l] = true
		next, err := os.Readlink(l)
		if err != nil {
			return
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

// runinterp runs the interpreter with the --list switch
// and the file as an argument. For each returned line
// it looks for => as the second field, indicating a
// real .so (as opposed to the .vdso or a string like
// 'not a dynamic executable'.
func runinterp(interp, file string) ([]string, error) {
	var names []string
	o, err := exec.Command(interp, "--list", file).Output()
	if err != nil {
		return nil, err
	}
	for _, p := range strings.Split(string(o), "\n") {
		f := strings.Split(p, " ")
		if len(f) < 3 {
			continue
		}
		if f[1] != "=>" || len(f[2]) == 0 {
			continue
		}
		names = append(names, f[2])
	}
	return names, nil
}

// Ldd returns a list of all library dependencies for a
// set of files, suitable for feeding into (e.g.) a cpio
// program. If a file has no dependencies, that is not an
// error. The only possible error is if a file does not
// exist, or it says it has an interpreter but we can't read
// it, or we are not able to run its interpreter.
// It's not an error for a file to not be an ELF, as
// this function should be convenient and the list might
// include non-ELF executables (a.out format, scripts)
func Ldd(names []string) ([]string, error) {
	var (
		list    = make(map[string]bool)
		interps = make(map[string]bool)
		libs    []string
	)
	for _, n := range names {
		r, err := os.Open(n)
		if err != nil {
			return nil, err
		}
		defer r.Close()
		f, err := elf.NewFile(r)
		if err != nil {
			continue
		}
		s := f.Section(".interp")
		if s == nil {
			continue
		}
		// If there is an interpreter, it should be
		// an error if we can't read it.
		i, err := s.Data()
		if err != nil {
			return nil, err
		}
		// Ignore #! interpreters
		if len(i) > 1 && i[0] == '#' && i[1] == '!' {
			continue
		}
		// annoyingly, s.Data() seems to return the null at the end and,
		// weirdly, that seems to confuse the kernel. Truncate it.
		interp := string(i[:len(i)-1])
		// We could just append the interp but people
		// expect to see that first.
		if !interps[interp] {
			interps[interp] = true
			libs = append(libs, interp)
		}
		// oh boy. Now to run the interp and get more names.
		n, err := runinterp(interp, n)
		if err != nil {
			return nil, err
		}
		for i := range n {
			follow(n[i], list)
		}
	}

	for i := range list {
		libs = append(libs, i)
	}

	return libs, nil
}
