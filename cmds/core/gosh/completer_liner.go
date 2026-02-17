// Copyright 2021 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!tinygo || tinygo.enable) && !plan9 && !goshsmall && goshliner

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// HistFile is the history file.
// This might, possibly, use GetPid to avoid gosh'es writing over each other
var HistFile = filepath.Join(os.TempDir(), "gosh.history")

var completion = flag.Bool("comp", true, "Enable tabcompletion and a more feature rich editline implementation")

func runInteractive(runner *interp.Runner, parser *syntax.Parser, stdout, stderr io.Writer) error {
	var hist io.Writer
	if f, err := os.OpenFile(HistFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o600); err == nil {
		defer f.Close()
		hist = f
	}
	_ = completion // completion is unsupported in minimal mode.
	// Ignore SIGINT in the shell itself; foreground child processes still receive it.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	defer signal.Stop(ch)
	go func(ch chan os.Signal) {
		for range ch {
		}
	}(ch)

	var runErr error
	for {
		if runErr != nil {
			fmt.Fprintf(stdout, "error: %s\n", runErr.Error())
			runErr = nil
		}

		line, err := readMinimalLine(os.Stdin, stdout, "$ ", parser, *completion)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			if !errors.Is(err, errPromptAborted) {
				fmt.Fprintf(stderr, "error: %s\n", err.Error())
			}
			continue
		}

		switch line {
		case "exit":
			goto exit
		case "disablecomp":
		case "enablecomp":
			continue
		default:
		}

		if line != "" {
			if hist != nil {
				fmt.Fprintf(hist, "%s\n", line)
			}
		}
		if err := parser.Stmts(strings.NewReader(line), func(stmt *syntax.Stmt) bool {
			if parser.Incomplete() {
				fmt.Fprintf(stdout, "-> ")
				return true
			}

			runErr = runner.Run(context.Background(), stmt)
			return !runner.Exited()
		}); err != nil {
			fmt.Fprintf(stderr, "error: %s\n", err.Error())
		}
	}
exit:
	return nil
}
