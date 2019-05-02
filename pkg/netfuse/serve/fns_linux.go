package serve

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"bazil.org/fuse"
)

type Timespec struct {
	Sec  int64
	Nsec int64
}

func spec(t syscall.Timespec) time.Time {
	return time.Unix(t.Sec, t.Nsec)
}

func Attr(n string, attr *fuse.Attr) error {
	fi, err := os.Stat(n)
	if err != nil {
		return err
	}
	attr.Valid = attrValidTime
	l, ok := fi.Sys().(*syscall.Stat_t)
	if !ok || l == nil {
		return fmt.Errorf("Can't stat without fi.Sys()")
	}
	attr.Inode = l.Ino
	attr.Size = uint64(l.Size)
	attr.Blocks = uint64(l.Blocks)
	attr.Atime = spec(l.Atim)
	attr.Mtime = spec(l.Mtim)
	attr.Ctime = spec(l.Ctim)
	attr.Mode = fi.Mode()
	Debug("Set mode to %#o", attr.Mode)
	attr.Nlink = uint32(l.Nlink)
	attr.Uid = l.Uid
	attr.Gid = l.Gid
	//attr.Dev = l.Dev
	attr.Rdev = uint32(l.Rdev)
	attr.Flags = 0x0
	attr.BlockSize = uint32(l.Blksize)
	return nil
}
