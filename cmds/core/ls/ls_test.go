// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/u-root/u-root/pkg/ls"
	"golang.org/x/sys/unix"
)

type ttflags struct {
	all       bool
	human     bool
	directory bool
	long      bool
	quoted    bool
	recurse   bool
	classify  bool
}

// Test listName func
func TestListName(t *testing.T) {
	// Create some directorys.
	tmpDir := t.TempDir()
	if err := os.Mkdir(filepath.Join(tmpDir, "d1"), 0777); err != nil {
		t.Fatalf("err in os.Mkdir: %v", err)
	}
	// Create some files.
	files := []string{"f1", "f2", "f3\nline 2", ".f4", "d1/f4"}
	for _, file := range files {
		if _, err := os.Create(filepath.Join(tmpDir, file)); err != nil {
			t.Fatalf("err in os.Create: %v", err)
		}
	}

	// Creating test table
	for _, tt := range []struct {
		name   string
		input  string
		want   string
		flag   ttflags
		prefix bool
	}{
		{
			name:  "ls without arguments",
			input: tmpDir,
			want:  fmt.Sprintf("%s\n%s\n%s\n%s\n", "d1", "f1", "f2", "f3?line 2"),
			flag:  ttflags{},
		},
		{
			name:  "ls osfi.IsDir() path, quoted = true",
			input: tmpDir,
			want:  fmt.Sprintf("\"%s\"\n\"%s\"\n\"%s\"\n\"%s\"\n", "d1", "f1", "f2", "f3\\nline 2"),
			flag: ttflags{
				quoted: true,
			},
		},
		{
			name:  "ls osfi.IsDir() path, quoted = true, prefix = true ",
			input: tmpDir,
			want:  fmt.Sprintf("\"%s\":\n\"%s\"\n\"%s\"\n\"%s\"\n\"%s\"\n", tmpDir, "d1", "f1", "f2", "f3\\nline 2"),
			flag: ttflags{
				quoted: true,
			},
			prefix: true,
		},
		{
			name:   "ls osfi.IsDir() path, quoted = false, prefix = true ",
			input:  tmpDir,
			want:   fmt.Sprintf("%s:\n%s\n%s\n%s\n%s\n", tmpDir, "d1", "f1", "f2", "f3?line 2"),
			prefix: true,
		},
		{
			name:  "ls recurse = true",
			input: tmpDir,
			want: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n", tmpDir, filepath.Join(tmpDir, ".f4"), filepath.Join(tmpDir, "d1"),
				filepath.Join(tmpDir, "d1/f4"), filepath.Join(tmpDir, "f1"), filepath.Join(tmpDir, "f2"), filepath.Join(tmpDir, "f3?line 2")),
			flag: ttflags{
				all:     true,
				recurse: true,
			},
		},
		{
			name:  "ls directory = true",
			input: tmpDir,
			want:  fmt.Sprintf("%s\n", "tmpDir"),
			flag: ttflags{
				directory: true,
			},
		},
		{
			name:  "ls classify = true",
			input: tmpDir,
			want:  fmt.Sprintf("%s\n%s\n%s\n%s\n", "d1/", "f1", "f2", "f3?line 2"),
			flag: ttflags{
				classify: true,
			},
		},
		{
			name:  "file does not exist",
			input: "dir",
			flag:  ttflags{},
		},
	} {
		// Setting the flags
		*all = tt.flag.all
		*human = tt.flag.human
		*directory = tt.flag.directory
		*long = tt.flag.long
		*quoted = tt.flag.quoted
		*recurse = tt.flag.recurse
		*classify = tt.flag.classify

		// Running the tests
		t.Run(tt.name, func(t *testing.T) {
			// Write output in buffer.
			var buf bytes.Buffer
			var s ls.Stringer = ls.NameStringer{}
			if *quoted {
				s = ls.QuotedStringer{}
			}
			if *long {
				s = ls.LongStringer{Human: *human, Name: s}
			}
			if err := listName(s, tt.input, &buf, tt.prefix); err != nil {
				if buf.String() != tt.want {
					t.Errorf("listName() = '%v', want: '%v'", buf.String(), tt.want)
				}
			}
		})
	}
}

// Test list func
func TestList(t *testing.T) {
	// Creating test table
	for _, tt := range []struct {
		name   string
		input  []string
		err    error
		flag   ttflags
		prefix bool
	}{
		{
			name:  "input empty, quoted = true, long = true",
			input: []string{},
			err:   nil,
			flag: ttflags{
				quoted: true,
				long:   true,
			},
		},
		{
			name:  "input empty, quoted = true, long = true",
			input: []string{"dir"},
			err:   os.ErrNotExist,
		},
	} {

		// Setting the flags
		*long = tt.flag.long
		*quoted = tt.flag.quoted

		// Running the tests
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := list(&buf, tt.input); err != nil {
				if !errors.Is(err, tt.err) {
					t.Errorf("list() = '%v', want: '%v'", err, tt.err)
				}
			}
		})
	}
}

// Test indicator func
func TestIndicator(t *testing.T) {
	// Creating test table
	for _, test := range []struct {
		lsInfo ls.FileInfo
		symbol string
	}{
		{
			ls.FileInfo{
				Mode: os.ModeDir,
			},
			"/",
		},
		{
			ls.FileInfo{
				Mode: os.ModeNamedPipe,
			},
			"|",
		},
		{
			ls.FileInfo{
				Mode: os.ModeSymlink,
			},
			"@",
		},
		{
			ls.FileInfo{
				Mode: os.ModeSocket,
			},
			"=",
		},
		{
			ls.FileInfo{
				Mode: 0b110110100,
			},
			"",
		},
		{
			ls.FileInfo{
				Mode: 0b111111101,
			},
			"*",
		},
	} {
		// Run tests
		got := indicator(test.lsInfo)
		if got != test.symbol {
			t.Errorf("for mode '%b' expected '%q', got '%q'", test.lsInfo.Mode, test.symbol, got)
		}
	}
}

// Make sure if perms fail in a dir, we still list the dir.
func TestPermHandling(t *testing.T) {
	d := t.TempDir()
	for _, v := range []string{"a", "c", "d"} {
		if err := os.Mkdir(filepath.Join(d, v), 0777); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(d, "b"), 0); err != nil {
		t.Fatal(err)
	}
	for _, v := range []string{"0", "1", "2"} {
		if err := os.Mkdir(filepath.Join(d, v), 0777); err != nil {
			t.Fatal(err)
		}
	}
	b := &bytes.Buffer{}
	if err := listName(ls.NameStringer{}, d, b, false); err != nil {
		t.Fatalf("listName(ls.NameString{}, %q, w, false): %v != nil", d, err)
	}
	// the output varies very widely between kernels and Go versions :-(
	// Just look for 'permission denied' and more than 6 lines of output ...
	if !strings.Contains(b.String(), "0\n1\n2\na\nb\nc\nd\n") {
		t.Errorf("ls %q: output %q did not contain %q", d, b.String(), "0\n1\n2\na\nb\nc\nd\n")
	}
}

func TestNotExist(t *testing.T) {
	d := t.TempDir()
	b := &bytes.Buffer{}
	if err := listName(ls.NameStringer{}, filepath.Join(d, "b"), b, false); err != nil {
		t.Fatalf("listName(ls.NameString{}, %q/b, w, false): nil != %v", d, err)
	}
	// yeesh.
	// errors not consistent and ... the error has this gratuitous 'lstat ' in front
	// of the filename ...
	eexist := fmt.Sprintf("%s:%v", filepath.Join(d, "b"), os.ErrNotExist)
	enoent := fmt.Sprintf("%s: %v", filepath.Join(d, "b"), unix.ENOENT)
	if !strings.Contains(b.String(), eexist) && !strings.Contains(b.String(), enoent) {
		t.Fatalf("ls of bad name: %q does not contain %q or %q", b.String(), eexist, enoent)
	}
}
