// FUSE service loop, for servers that wish to use it.

package serve

import (
	"os"
	"time"

	"bazil.org/fuse"
	"golang.org/x/sys/unix"
)

type FileSystem struct {
	root string
}

func NewFileSystem(root string) (FS, error) {
	Debug("NewFileSystem: root %q", root)
	return &FileSystem{
		root: root,
	}, nil
}

func (f *FileSystem) Getattr(req *fuse.GetattrRequest, resp *fuse.GetattrResponse) error {
	Debug("Attr: %v %v", req, nil)
	i, err := os.Stat(f.root)
	if err != nil {
		return err
	}
	Debug("Stat %v get %v", f.root, i)
	/*
		resp.Valid = attrValidTime
		resp.Size = uint64(i.Size())
		resp.Blocks = resp.Size / 4096
		resp.Mode = i.Mode()
	*/
	return nil
	//type FileInfo interface {
	// Name() string       // base name of the file
	// Size() int64        // length in bytes for regular files; system-dependent for others
	// Mode() FileMode     // file mode bits
	// ModTime() time.Time // modification time
	// IsDir() bool        // abbreviation for Mode().IsDir()
	// Sys() interface{}   // underlying data source (can return nil)

}

func (f *FileSystem) Statfs(req *fuse.StatfsRequest, resp *fuse.StatfsResponse) error {
	var buf unix.Statfs_t
	err := unix.Statfs(f.root, &buf)
	if err != nil {
		return err
	}
	var _ = &fuse.Attr{
		//Inode:  buf.Fsid,
		Size:   uint64(buf.Bsize) * buf.Blocks,
		Blocks: uint64(buf.Bsize) * buf.Blocks / 512,
		Atime:  time.Now(),
		Mtime:  time.Now(),
		Ctime:  time.Now(),
		/*		Mode
				Mode      os.FileMode // file mode
				Nlink     uint32      // number of links (usually 1)
				Uid       uint32      // owner uid
				Gid       uint32      // group gid
				Rdev      uint32      // device numbers
				Flags     uint32      // chflags(2) flags (OS X only)
				BlockSize uint32      // preferred blocksize for filesystem I/O
		*/

		/*
			type Statfs_t struct {
				Type    int64
				Bsize   int64
				Blocks  uint64
				Bfree   uint64
				Bavail  uint64
				Files   uint64
				Ffree   uint64
				Fsid    Fsid
				Namelen int64
				Frsize  int64
				Flags   int64
				Spare   [4]int64
			}*/

		Valid: attrValidTime,

		/*		Mode
				Mode      os.FileMode // file mode
				Nlink     uint32      // number of links (usually 1)
				Uid       uint32      // owner uid
				Gid       uint32      // group gid
				Rdev      uint32      // device numbers
				Flags     uint32      // chflags(2) flags (OS X only)
				BlockSize uint32      // preferred blocksize for filesystem I/O
		*/
	}

	return nil
}
