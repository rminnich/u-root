// pstore implements a line-oriented persistent store via Fuse.
// reads return the contents of the file.
// Write to the file are only committed once the file is closed.
// They are always appended (i.e. offset is ignored)
// one use of the pstore is for history files.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"
        _ "bazil.org/fuse/fs/fstestutil"
	"golang.org/x/net/context"
)

var (
	store *os.File
	rw sync.RWMutex
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
	var err error
	store, err = os.OpenFile(*backingStore, os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("%v", err)
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
		return File{d: &bytes.Buffer{}}, nil
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
type File struct{d *bytes.Buffer}

func (File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 2
	a.Mode = 0644
	end, err := store.Seek(0, 2)
	if err != nil {
		return err
	}
	a.Size = uint64(end)
	return nil
}

func (f File) ReadAll(ctx context.Context) ([]byte, error) {
	rw.RLock()
	defer rw.RUnlock()
	var a fuse.Attr
	if err := f.Attr(ctx, &a); err != nil {
		return nil, err
	}
	d := make([]byte, a.Size)
	_, err := store.ReadAt(d, 0)
	return d, err
}

var _ fs.NodeOpener = (*File)(nil)

func (f File) Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (fs.Handle, error) {
	resp.Flags |= fuse.OpenKeepCache
	return f, nil
}

var _ fs.Handle = (*File)(nil)

var _ fs.HandleReader = (*File)(nil)

func (f File) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	rw.RLock()
	defer rw.RUnlock()
	var data = make([]byte, req.Size)
	_, err := store.Read(data)
	if err != nil {
		data = nil
	}
	fuseutil.HandleRead(req, resp, nil)
	return err
}

var _ fs.HandleWriter = (*File)(nil)

func (f File) Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	f.d.Write(req.Data)
	return nil
}

var _ fs.HandleReleaser = (*File)(nil)

func (f File) Release(ctx context.Context, req *fuse.ReleaseRequest) error {
	rw.Lock()
	var a fuse.Attr
	if err := f.Attr(ctx, &a); err != nil {
		return err
	}
	defer rw.Unlock()
	_, err := store.WriteAt(f.d.Bytes(), int64(a.Size))
	return err
}

