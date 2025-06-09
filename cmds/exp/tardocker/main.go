package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
)

type cmd struct {
	container string
}

func (c *cmd) Run() error {
	ref, err := name.ParseReference(c.container)
	if err != nil {
		return err
	}

	img, err := crane.Pull(ref.Name())
	if err != nil {
		return err
	}

	r := mutate.Extract(img)

	if _, err := io.Copy(os.Stdout, r); err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Fatalf("Usage: %s <image>", os.Args[0])
	}
	c := &cmd{container: flag.Arg(0)}
	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
}
