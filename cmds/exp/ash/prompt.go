// Copyright 2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"
)

var (
	ps = "u@h:d% "
	home = os.Getenv("HOME")
	user = os.Getenv("USER")
	host = os.Getenv("HOST")

)

func init() {
	addBuiltIn("prompt", promptBuiltin)
	if host == "" {
		host = "unknown"
		if h, err := os.Hostname(); err == nil {
			host = h
		}
	}
}

func prompt() string {
	var out, v string
	var err error
	for _, c := range ps {
		switch c {
		case 'd':
			v, err = os.Getwd()
			if strings.HasPrefix(v, home) {
				v = "~" + v[len(home):]
			}
		case 'h':
			v = host
		case 'u':
			v = user
		default:
			v, err = string(c), nil
		}
		if err != nil {
			out = out + fmt.Sprintf("%v", err)
		}
		if v != "" {
			out += v
		}
	}
	return out
}

func promptBuiltin(c *Command) error {
	if len(c.argv) != 1 {
		return fmt.Errorf("Usage: prompt format")
	}
	ps = c.argv[0]
	return nil
}
