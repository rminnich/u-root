// Hellofs implements a simple "hello world" file system.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"
)

var (
	data []byte
	backingStore = flag.String("store", "store", "name of file holding all the data")
	fileName = flag.String("filename", "data", "name of file in file system")
	fsName = flag.String("fsname", "pstore", "Name of file system")
	subType = flag.String("subtype", "pstorefs", "subtype of file system")
	volumeName = flag.String("volumeName", "Persistent Store", "volumeName of file system")
)


func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s MOUNTPOINT\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		os.Exit(2)
	}
	mountpoint := flag.Arg(0)

	// First, make sure we can operate on the data
	f, err := os.OpenFile(*backingStore, os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("%v", err)
	}

	data, err = ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Reading %v: %v", f, err)
	}

	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName(*fsName),
		fuse.Subtype(*subType),
		fuse.LocalVolume(),
		fuse.VolumeName(*volumeName))
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	err = fs.Serve(c, FS{})
	if err != nil {
		log.Fatal(err)
	}

	// check if the mount process has an error to report
	<-c.Ready
	if err := c.MountError; err != nil {
		log.Fatal(err)
	}
}

// FS implements the hello world file system.
type FS struct{}

func (FS) Root() (fs.Node, error) {
	return Dir{}, nil
}

// Dir implements both Node and Handle for the root directory.
type Dir struct{}

func (Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Mode = os.ModeDir | 0555
	return nil
}

func (Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if name == *fileName {
		return File{}, nil
	}
	return nil, fuse.ENOENT
}

var dirDirs = []fuse.Dirent{
	{Inode: 2, Name: *fileName, Type: fuse.DT_File},
}

func (Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	return dirDirs, nil
}

// File implements both Node and Handle for the hello file.
type File struct{}

func (File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 2
	a.Mode = 0444
	a.Size = uint64(len(data))
	return nil
}

func (File) ReadAll(ctx context.Context) ([]byte, error) {
	return data, nil
}
