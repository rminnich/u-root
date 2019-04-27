// FUSE service loop, for servers that wish to use it.

package serve

import (
	"bazil.org/fuse"
	"golang.org/x/net/context"
)

type FileSystem struct {
	root string
}

func NewFileSystem(root string) (FS, error) {
	return &FileSystem{
		root: root,
	}, nil
}

func (f *FileSystem) Root() (Node, error) {
	return f, nil
}

func (f *FileSystem) Attr(ctx context.Context, attr *fuse.Attr) error {
	Debug("Attr: %v %v", attr, nil)
	return nil
}
