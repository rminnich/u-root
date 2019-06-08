// Copyright 2012-2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/u-root/u-root/pkg/complete"
	"github.com/u-root/u-root/pkg/termios"
)

var (
	debug   = flag.Bool("d", false, "enable debug prints")
	userush = flag.Bool("R", true, "Use the rush interpreter for commands")
	v       = func(string, ...interface{}) {}
)

func verbose(f string, a ...interface{}) {
	v(f+"\r\n", a...)
}

func output(s chan string, w io.Writer) {
	var prevline string
	for l := range s {
		p := "\r"
		if l == "\n" {
			fmt.Print("\r\n"+prompt())
			prevline = prompt()
			continue
		}
		if _, err := w.Write([]byte(p)); err != nil {
			log.Printf("output write: %v", err)
		}
		// We are now at the left. Write blanks for the size of prevline
		// followed by newline
		if len(prevline) > len(prompt()) + len(l) {
			fmt.Fprintf(w, "%*s", len(prevline)+1, "\r")
		}
		fmt.Print(prompt())
		prevline = prompt()
		for _, b := range l {
			var o string
			switch b {
			default:
				o = string(b)
				prevline += o
			case '\r', '\n':
				o = "\r\n"
				prevline = ""
			}
			if _, err := w.Write([]byte(o)); err != nil {
				log.Printf("output write: %v", err)
			}
		}
	}
}

func main() {
	tty()
	flag.Parse()
	if *debug {
		v = log.Printf
		complete.Debug = verbose
	}
	t, err := termios.New()
	if err != nil {
		log.Fatal(err)
	}
	r, err := t.Raw()
	if err != nil {
		log.Printf("non-fatal cannot get tty: %v", err)
	}
	defer t.Set(r)
	_, cw, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}
	p, err := complete.NewPathCompleter()
	if err != nil {
		log.Fatal(err)
	}
	f := complete.NewFileCompleter("/")
	bin := complete.NewMultiCompleter(complete.NewStringCompleter(builtinList), p)
	rest := f
	l := complete.NewLineReader(bin, t, cw)
	lines := make(chan string)
	go output(lines, os.Stdout)
	var lineComplete bool
	fmt.Print(prompt())
	for !l.EOF {
		lineComplete = false
		l.C = bin
		if l.Fields > 1 {
			l.C = rest
		}
		// Read one byte, run it through the completer, then print the string
		// as we have it.
		v("start with %v", l)
		var b [1]byte
		n, err := l.R.Read(b[:])
		if err != nil {
			break
		}
		v("ReadLine: got %s, %v, %v", b, n, err)

		if err := l.ReadChar(b[0]); err != nil {
			v("ERR -> %v (%v)", l, err)
			if err == io.EOF || err != complete.ErrEOL {
				v("%v", err)
				lines <- l.Line
				continue
			}
			v("set linecomplete")
			lineComplete = true
		}

		v("back from ReadChar, l is %v", l)
		if lineComplete && l.Line != "" {
			v("ash: Done reading args: line %q", l.Line)
			// here we go.
			lines <- l.Line+"\r"
			t.Set(r)
			if !*userush {
				f := strings.Fields(l.Line)
				var args []string
				if l.Exact != "" {
					args = append(args, l.Exact)
				}
				args = append(args, l.Candidates...)
				if len(f) > 1 && len(args) > 1 {
					f = append(f[:len(f)-1], args...)
				}

				cmd := exec.Command(f[0], f[1:]...)
				cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

				if err := cmd.Run(); err != nil {
					log.Print(err)
				}

			} else {
				if l.Line == "exit" {
					break
				}
				b := bufio.NewReader(bytes.NewBufferString(l.Line))
				if err := rush(b); err != nil {
					log.Print(err)
				}
			}
			foreground()
			t.Raw()
			lines <- "\n"
			l.Line = ""
			l.Candidates = []string{}
			l.C = bin
			l.Fields = 0
			l.Exact = ""
			continue
		}
		if lineComplete {
			lines <- "\n"
			l.Line = ""
			l.Candidates = []string{}
			l.C = bin
			l.Fields = 0
			l.Exact = ""
			continue
		}
		if l.Exact != "" {
			lines <- "\n" + l.Exact
		}
		if len(l.Candidates) > 0 {
			for _, ln := range l.Candidates {
				lines <- "\n" + ln
			}
			lines <- strings.Join(l.Candidates, "\n") + "\n"
		}
		lines <- l.Line
	}
}
