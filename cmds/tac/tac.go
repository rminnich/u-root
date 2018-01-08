package main

import (
	"flag"
	"io"
	"log"
	"os"
	"sync"
)

const ReadSize int64 = 4096

type ReadAtSeeker interface {
	io.ReaderAt
	io.Seeker
}

func tac(r ReadAtSeeker, w io.Writer) {
	var b [ReadSize]byte
	// Get current EOF. While the file may be growing, there's
	// only so much we can do.
	loc, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	c := make(chan byte)
	go func(r <-chan byte, w io.Writer) {
		defer wg.Done()
		line := string(<-r)
		for c := range r {
			if c == '\n' {
				if _, err := w.Write([]byte(line)); err != nil {
					log.Fatal(err)
				}
				line = ""
			}
			line = string(c) + line
		}
		if _, err := w.Write([]byte(line)); err != nil {
			log.Fatal(err)
		}
	}(c, w)

	for loc > 0 {
		n := ReadSize
		if loc < ReadSize {
			n = loc
		}

		amt, err := r.ReadAt(b[:n], loc-int64(n))
		if err != nil {
			log.Printf("%v ", err)
			break
		}
		loc -= int64(amt)
		for i := range b[:amt] {
			o := amt - i
			c <- b[o]
		}
	}
	close(c)
	wg.Wait()
}

func main() {
	flag.Parse()
	a := flag.Args()
	if len(a) != 1 {
		log.Fatalf("no")
	}

	f, err := os.Open(a[0])
	if err != nil {
		log.Fatal(err)
	}
	tac(f, os.Stdout)
}
