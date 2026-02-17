// Copyright 2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Minimal interactive line reader without ANSI escape sequences.

//go:build (!tinygo || tinygo.enable) && !plan9 && !goshsmall

package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
	"mvdan.cc/sh/v3/syntax"
)

var errPromptAborted = errors.New("prompt aborted")

func readMinimalLine(in *os.File, out io.Writer, prompt string, parser *syntax.Parser, enableCompletion bool) (string, error) {
	if _, err := fmt.Fprint(out, prompt); err != nil {
		return "", err
	}

	oldState, err := term.MakeRaw(int(in.Fd()))
	if err != nil {
		return "", err
	}
	defer term.Restore(int(in.Fd()), oldState)

	buf := make([]byte, 0, 128)
	tmp := make([]byte, 1)

	for {
		_, err := in.Read(tmp)
		if err != nil {
			return "", err
		}
		b := tmp[0]

		switch b {
		case '\r', '\n':
			fmt.Fprint(out, "\r\n")
			return string(buf), nil
		case 3: // Ctrl-C
			fmt.Fprint(out, "^C\r\n")
			return "", errPromptAborted
		case 4: // Ctrl-D
			if len(buf) == 0 {
				return "", io.EOF
			}
		case 8, 127: // backspace/delete
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
				fmt.Fprint(out, "\b \b")
			}
		case '\t':
			if !enableCompletion {
				continue
			}

			line := string(buf)
			matches := autocompleteLiner(parser)(line)
			switch len(matches) {
			case 0:
			case 1:
				m := matches[0]
				if strings.HasPrefix(m, line) && len(m) > len(line) {
					suffix := m[len(line):]
					buf = append(buf, []byte(suffix)...)
					fmt.Fprint(out, suffix)
				} else if m != line {
					// For normalized/expanded matches (e.g. quoted args), rewrite the full line.
					for range len(buf) {
						fmt.Fprint(out, "\b \b")
					}
					buf = append(buf[:0], []byte(m)...)
					fmt.Fprint(out, m)
				}
			default:
				fmt.Fprint(out, "\r\n")
				for _, m := range matches {
					fmt.Fprintf(out, "\r%s\r\n", m)
				}
				fmt.Fprintf(out, "\r%s%s", prompt, line)
			}
		default:
			// Keep the input path simple and byte-oriented in minimal mode.
			if b >= 32 {
				buf = append(buf, b)
				fmt.Fprintf(out, "%c", b)
			}
		}
	}
}
