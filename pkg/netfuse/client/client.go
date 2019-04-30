// FUSE service loop, for servers that wish to use it.

package serve

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"

	"bytes"

	"bazil.org/fuse"
	"golang.org/x/net/context"
)

var Debug = log.Printf

// TODO: FINISH DOCS

// An FS is the interface required of a file system.
//
// Other FUSE requests can be handled by implementing methods from the
// FS* interfaces, for example FSStatfser.
type FS interface {
	// Root is called to obtain the Node for the file system root.
	Root() (Node, error)
}

type FSStatfser interface {
	// Statfs is called to obtain file system metadata.
	// It should write that data to resp.
	Statfs(ctx context.Context, req *fuse.StatfsRequest, resp *fuse.StatfsResponse) error
}

type FSDestroyer interface {
	// Destroy is called when the file system is shutting down.
	//
	// Linux only sends this request for block device backed (fuseblk)
	// filesystems, to allow them to flush writes to disk before the
	// unmount completes.
	Destroy()
}

type FSInodeGenerator interface {
	// GenerateInode is called to pick a dynamic inode number when it
	// would otherwise be 0.
	//
	// Not all filesystems bother tracking inodes, but FUSE requires
	// the inode to be set, and fewer duplicates in general makes UNIX
	// tools work better.
	//
	// Operations where the nodes may return 0 inodes include Getattr,
	// Setattr and ReadDir.
	//
	// If FS does not implement FSInodeGenerator, GenerateDynamicInode
	// is used.
	//
	// Implementing this is useful to e.g. constrain the range of
	// inode values used for dynamic inodes.
	GenerateInode(parentInode uint64, name string) uint64
}

// A Node is the interface required of a file or directory.
// See the documentation for type FS for general information
// pertaining to all methods.
//
// A Node must be usable as a map key, that is, it cannot be a
// function, map or slice.
//
// Other FUSE requests can be handled by implementing methods from the
// Node* interfaces, for example NodeOpener.
//
// Methods returning Node should take care to return the same Node
// when the result is logically the same instance. Without this, each
// Node will get a new NodeID, causing spurious cache invalidations,
// extra lookups and aliasing anomalies. This may not matter for a
// simple, read-only filesystem.
type Node interface {
	// Attr fills attr with the standard metadata for the node.
	//
	// Fields with reasonable defaults are prepopulated. For example,
	// all times are set to a fixed moment when the program started.
	//
	// If Inode is left as 0, a dynamic inode number is chosen.
	//
	// The result may be cached for the duration set in Valid.
	Attr(ctx context.Context, attr *fuse.Attr) error
}

type NodeGetattrer interface {
	// Getattr obtains the standard metadata for the receiver.
	// It should store that metadata in resp.
	//
	// If this method is not implemented, the attributes will be
	// generated based on Attr(), with zero values filled in.
	Getattr(ctx context.Context, req *fuse.GetattrRequest, resp *fuse.GetattrResponse) error
}

type NodeSetattrer interface {
	// Setattr sets the standard metadata for the receiver.
	//
	// Note, this is also used to communicate changes in the size of
	// the file, outside of Writes.
	//
	// req.Valid is a bitmask of what fields are actually being set.
	// For example, the method should not change the mode of the file
	// unless req.Valid.Mode() is true.
	Setattr(ctx context.Context, req *fuse.SetattrRequest, resp *fuse.SetattrResponse) error
}

type NodeSymlinker interface {
	// Symlink creates a new symbolic link in the receiver, which must be a directory.
	//
	// TODO is the above true about directories?
	Symlink(ctx context.Context, req *fuse.SymlinkRequest) (Node, error)
}

// This optional request will be called only for symbolic link nodes.
type NodeReadlinker interface {
	// Readlink reads a symbolic link.
	Readlink(ctx context.Context, req *fuse.ReadlinkRequest) (string, error)
}

type NodeLinker interface {
	// Link creates a new directory entry in the receiver based on an
	// existing Node. Receiver must be a directory.
	Link(ctx context.Context, req *fuse.LinkRequest, old Node) (Node, error)
}

type NodeRemover interface {
	// Remove removes the entry with the given name from
	// the receiver, which must be a directory.  The entry to be removed
	// may correspond to a file (unlink) or to a directory (rmdir).
	Remove(ctx context.Context, req *fuse.RemoveRequest) error
}

type NodeAccesser interface {
	// Access checks whether the calling context has permission for
	// the given operations on the receiver. If so, Access should
	// return nil. If not, Access should return EPERM.
	//
	// Note that this call affects the result of the access(2) system
	// call but not the open(2) system call. If Access is not
	// implemented, the Node behaves as if it always returns nil
	// (permission granted), relying on checks in Open instead.
	Access(ctx context.Context, req *fuse.AccessRequest) error
}

type NodeStringLookuper interface {
	// Lookup looks up a specific entry in the receiver,
	// which must be a directory.  Lookup should return a Node
	// corresponding to the entry.  If the name does not exist in
	// the directory, Lookup should return ENOENT.
	//
	// Lookup need not to handle the names "." and "..".
	Lookup(ctx context.Context, name string) (Node, error)
}

type NodeRequestLookuper interface {
	// Lookup looks up a specific entry in the receiver.
	// See NodeStringLookuper for more.
	Lookup(ctx context.Context, req *fuse.LookupRequest, resp *fuse.LookupResponse) (Node, error)
}

type NodeMkdirer interface {
	Mkdir(ctx context.Context, req *fuse.MkdirRequest) (Node, error)
}

type NodeOpener interface {
	// Open opens the receiver. After a successful open, a client
	// process has a file descriptor referring to this Handle.
	//
	// Open can also be also called on non-files. For example,
	// directories are Opened for ReadDir or fchdir(2).
	//
	// If this method is not implemented, the open will always
	// succeed, and the Node itself will be used as the Handle.
	//
	// XXX note about access.  XXX OpenFlags.
	Open(ctx context.Context, req *fuse.OpenRequest, resp *fuse.OpenResponse) (Handle, error)
}

type NodeCreater interface {
	// Create creates a new directory entry in the receiver, which
	// must be a directory.
	Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (Node, Handle, error)
}

type NodeForgetter interface {
	// Forget about this node. This node will not receive further
	// method calls.
	//
	// Forget is not necessarily seen on unmount, as all nodes are
	// implicitly forgotten as part part of the unmount.
	Forget()
}

type NodeRenamer interface {
	Rename(ctx context.Context, req *fuse.RenameRequest, newDir Node) error
}

type NodeMknoder interface {
	Mknod(ctx context.Context, req *fuse.MknodRequest) (Node, error)
}

// TODO this should be on Handle not Node
type NodeFsyncer interface {
	Fsync(ctx context.Context, req *fuse.FsyncRequest) error
}

type NodeGetxattrer interface {
	// Getxattr gets an extended attribute by the given name from the
	// node.
	//
	// If there is no xattr by that name, returns fuse.ErrNoXattr.
	Getxattr(ctx context.Context, req *fuse.GetxattrRequest, resp *fuse.GetxattrResponse) error
}

type NodeListxattrer interface {
	// Listxattr lists the extended attributes recorded for the node.
	Listxattr(ctx context.Context, req *fuse.ListxattrRequest, resp *fuse.ListxattrResponse) error
}

type NodeSetxattrer interface {
	// Setxattr sets an extended attribute with the given name and
	// value for the node.
	Setxattr(ctx context.Context, req *fuse.SetxattrRequest) error
}

type NodeRemovexattrer interface {
	// Removexattr removes an extended attribute for the name.
	//
	// If there is no xattr by that name, returns fuse.ErrNoXattr.
	Removexattr(ctx context.Context, req *fuse.RemovexattrRequest) error
}

var startTime = time.Now()

func nodeAttr(ctx context.Context, n Node, attr *fuse.Attr) error {
	//	attr.Valid = attrValidTime
	attr.Nlink = 1
	attr.Atime = startTime
	attr.Mtime = startTime
	attr.Ctime = startTime
	attr.Crtime = startTime
	if err := n.Attr(ctx, attr); err != nil {
		return err
	}
	return nil
}

// A Handle is the interface required of an opened file or directory.
// See the documentation for type FS for general information
// pertaining to all methods.
//
// Other FUSE requests can be handled by implementing methods from the
// Handle* interfaces. The most common to implement are HandleReader,
// HandleReadDirer, and HandleWriter.
//
// TODO implement methods: Getlk, Setlk, Setlkw
type Handle interface {
}

type HandleFlusher interface {
	// Flush is called each time the file or directory is closed.
	// Because there can be multiple file descriptors referring to a
	// single opened file, Flush can be called multiple times.
	Flush(ctx context.Context, req *fuse.FlushRequest) error
}

type HandleReadAller interface {
	ReadAll(ctx context.Context) ([]byte, error)
}

type HandleReadDirAller interface {
	ReadDirAll(ctx context.Context) ([]fuse.Dirent, error)
}

type HandleReader interface {
	// Read requests to read data from the handle.
	//
	// There is a page cache in the kernel that normally submits only
	// page-aligned reads spanning one or more pages. However, you
	// should not rely on this. To see individual requests as
	// submitted by the file system clients, set OpenDirectIO.
	//
	// Note that reads beyond the size of the file as reported by Attr
	// are not even attempted (except in OpenDirectIO mode).
	Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error
}

type HandleWriter interface {
	// Write requests to write data into the handle at the given offset.
	// Store the amount of data written in resp.Size.
	//
	// There is a writeback page cache in the kernel that normally submits
	// only page-aligned writes spanning one or more pages. However,
	// you should not rely on this. To see individual requests as
	// submitted by the file system clients, set OpenDirectIO.
	//
	// Writes that grow the file are expected to update the file size
	// (as seen through Attr). Note that file size changes are
	// communicated also through Setattr.
	Write(ctx context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error
}

type HandleReleaser interface {
	Release(ctx context.Context, req *fuse.ReleaseRequest) error
}

type Config struct {
	// Function to send debug log messages to. If nil, use fuse.Debug.
	// Note that changing this or fuse.Debug may not affect existing
	// calls to Serve.
	//
	// See fuse.Debug for the rules that log functions must follow.
	Debug func(msg interface{})

	// Function to put things into context for processing the request.
	// The returned context must have ctx as its parent.
	//
	// Note that changing this may not affect existing calls to Serve.
	//
	// Must not retain req.
	WithContext func(ctx context.Context, req fuse.Request) context.Context
}

// Serve serves the FUSE connection by making calls to the methods
// of fs and the Nodes and Handles it makes available.  It returns only
// when the connection has been closed or an unexpected error occurs.
// func (s *Server) Serve(fs FS) error {
// 	defer s.wg.Wait() // Wait for worker goroutines to complete before return

// 	s.fs = fs
// 	if dyn, ok := fs.(FSInodeGenerator); ok {
// 		s.dynamicInode = dyn.GenerateInode
// 	}

// 	root, err := fs.Root()
// 	if err != nil {
// 		return fmt.Errorf("cannot obtain root node: %v", err)
// 	}
// 	// Recognize the root node if it's ever returned from Lookup,
// 	// passed to Invalidate, etc.
// 	s.nodeRef[root] = 1
// 	s.node = append(s.node, nil, &serveNode{
// 		inode:      1,
// 		generation: s.nodeGen,
// 		node:       root,
// 		refs:       1,
// 	})
// 	s.handle = append(s.handle, nil)

// 	for {
// 		req, err := s.conn.ReadRequest()
// 		if err != nil {
// 			if err == io.EOF {
// 				break
// 			}
// 			return err
// 		}

// 		s.wg.Add(1)
// 		go func() {
// 			defer s.wg.Done()
// 			s.serve(req)
// 		}()
// 	}
// 	return nil
// }

// // Serve serves a FUSE connection with the default settings. See
// // Server.Serve.
// func Serve(c *net.Conn, fs FS) error {
// 	server := New(c, nil)
// 	return server.Serve(fs)
// }

type nothing struct{}

type serveRequest struct {
	Request fuse.Request
	cancel  func()
}

type serveNode struct {
	inode      uint64
	generation uint64
	node       Node
	refs       uint64

	// Delay freeing the NodeID until waitgroup is done. This allows
	// using the NodeID for short periods of time without holding the
	// Server.meta lock.
	//
	// Rules:
	//
	//     - hold Server.meta while calling wg.Add, then unlock
	//     - do NOT try to reacquire Server.meta
	wg sync.WaitGroup
}

func (sn *serveNode) attr(ctx context.Context, attr *fuse.Attr) error {
	err := nodeAttr(ctx, sn.node, attr)
	if attr.Inode == 0 {
		attr.Inode = sn.inode
	}
	return err
}

type serveHandle struct {
	handle   Handle
	readData []byte
	nodeID   fuse.NodeID
}

type request struct {
	Op      string
	Request *fuse.Header
	In      interface{} `json:",omitempty"`
}

func (r request) String() string {
	return fmt.Sprintf("<- %s", r.In)
}

type logResponseHeader struct {
	ID fuse.RequestID
}

func (m logResponseHeader) String() string {
	return fmt.Sprintf("ID=%v", m.ID)
}

type response struct {
	Op      string
	Request logResponseHeader
	Out     interface{} `json:",omitempty"`
	// Errno contains the errno value as a string, for example "EPERM".
	Errno string `json:",omitempty"`
	// Error may contain a free form error message.
	Error string `json:",omitempty"`
}

func (r response) errstr() string {
	s := r.Errno
	if r.Error != "" {
		// prefix the errno constant to the long form message
		s = s + ": " + r.Error
	}
	return s
}

func (r response) String() string {
	switch {
	case r.Errno != "" && r.Out != nil:
		return fmt.Sprintf("-> [%v] %v error=%s", r.Request, r.Out, r.errstr())
	case r.Errno != "":
		return fmt.Sprintf("-> [%v] %s error=%s", r.Request, r.Op, r.errstr())
	case r.Out != nil:
		// make sure (seemingly) empty values are readable
		switch r.Out.(type) {
		case string:
			return fmt.Sprintf("-> [%v] %s %q", r.Request, r.Op, r.Out)
		case []byte:
			return fmt.Sprintf("-> [%v] %s [% x]", r.Request, r.Op, r.Out)
		default:
			return fmt.Sprintf("-> [%v] %v", r.Request, r.Out)
		}
	default:
		return fmt.Sprintf("-> [%v] %s", r.Request, r.Op)
	}
}

type notification struct {
	Op   string
	Node fuse.NodeID
	Out  interface{} `json:",omitempty"`
	Err  string      `json:",omitempty"`
}

func (n notification) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "=> %s %v", n.Op, n.Node)
	if n.Out != nil {
		// make sure (seemingly) empty values are readable
		switch n.Out.(type) {
		case string:
			fmt.Fprintf(&buf, " %q", n.Out)
		case []byte:
			fmt.Fprintf(&buf, " [% x]", n.Out)
		default:
			fmt.Fprintf(&buf, " %s", n.Out)
		}
	}
	if n.Err != "" {
		fmt.Fprintf(&buf, " Err:%v", n.Err)
	}
	return buf.String()
}

type logMissingNode struct {
	MaxNode fuse.NodeID
}

func opName(req fuse.Request) string {
	t := reflect.Indirect(reflect.ValueOf(req)).Type()
	s := t.Name()
	s = strings.TrimSuffix(s, "Request")
	return s
}

type logLinkRequestOldNodeNotFound struct {
	Request *fuse.Header
	In      *fuse.LinkRequest
}

func (m *logLinkRequestOldNodeNotFound) String() string {
	return fmt.Sprintf("In LinkRequest (request %v), node %d not found", m.Request.Hdr().ID, m.In.OldNode)
}

type renameNewDirNodeNotFound struct {
	Request *fuse.Header
	In      *fuse.RenameRequest
}

func (m *renameNewDirNodeNotFound) String() string {
	return fmt.Sprintf("In RenameRequest (request %v), node %d not found", m.Request.Hdr().ID, m.In.NewDir)
}

type handlerPanickedError struct {
	Request interface{}
	Err     interface{}
}

var _ error = handlerPanickedError{}

func (h handlerPanickedError) Error() string {
	return fmt.Sprintf("handler panicked: %v", h.Err)
}

var _ fuse.ErrorNumber = handlerPanickedError{}

func (h handlerPanickedError) Errno() fuse.Errno {
	if err, ok := h.Err.(fuse.ErrorNumber); ok {
		return err.Errno()
	}
	return fuse.DefaultErrno
}

// handlerTerminatedError happens when a handler terminates itself
// with runtime.Goexit. This is most commonly because of incorrect use
// of testing.TB.FailNow, typically via t.Fatal.
type handlerTerminatedError struct {
	Request interface{}
}

var _ error = handlerTerminatedError{}

func (h handlerTerminatedError) Error() string {
	return fmt.Sprintf("handler terminated (called runtime.Goexit)")
}

var _ fuse.ErrorNumber = handlerTerminatedError{}

func (h handlerTerminatedError) Errno() fuse.Errno {
	return fuse.DefaultErrno
}

type handleNotReaderError struct {
	handle Handle
}

var _ error = handleNotReaderError{}

func (e handleNotReaderError) Error() string {
	return fmt.Sprintf("handle has no Read: %T", e.handle)
}

var _ fuse.ErrorNumber = handleNotReaderError{}

func (e handleNotReaderError) Errno() fuse.Errno {
	return fuse.ENOTSUP
}

func initLookupResponse(s *fuse.LookupResponse) {
	//	s.EntryValid = entryValidTime
}

type Server struct {
	fuse.Conn
	Debug   func(msg interface{})
	Context func(ctx context.Context, req fuse.Request) context.Context

	// state, protected by meta
	meta       sync.Mutex
	req        map[fuse.RequestID]*serveRequest
	node       []*serveNode
	nodeRef    map[Node]fuse.NodeID
	handle     []*serveHandle
	freeNode   []fuse.NodeID
	freeHandle []fuse.HandleID
	nodeGen    uint64

	// Used to ensure worker goroutines finish before Serve returns
	wg sync.WaitGroup
}

func (c *Client) serve(r fuse.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	parentCtx := ctx
	if c.Context != nil {
		ctx = c.Context(ctx, r)
	}

	//req := &serveRequest{Request: r, cancel: cancel}

	c.Debug(request{
		Op:      opName(r),
		Request: r.Hdr(),
		In:      r,
	})

	// Call this before responding.
	// After responding is too late: we might get another request
	// with the same ID and be very confused.
	done := func(resp interface{}) {
		msg := response{
			Op:      opName(r),
			Request: logResponseHeader{ /*ID: hdr.ID*/ },
		}
		if err, ok := resp.(error); ok {
			msg.Error = err.Error()
			if ferr, ok := err.(fuse.ErrorNumber); ok {
				errno := ferr.Errno()
				msg.Errno = errno.ErrnoName()
				if errno == err {
					// it's just a fuse.Errno with no extra detail;
					// skip the textual message for log readability
					msg.Error = ""
				}
			} else {
				msg.Errno = fuse.DefaultErrno.ErrnoName()
			}
		} else {
			msg.Out = resp
		}
		c.Debug(msg)

		/*
			c.meta.Lock()
			delete(c.req, hdr.ID)
			c.meta.Unlock()*/
	}

	var responded bool
	defer func() {
		if rec := recover(); rec != nil {
			const size = 1 << 16
			buf := make([]byte, size)
			n := runtime.Stack(buf, false)
			buf = buf[:n]
			log.Printf("fuse: panic in handler for %v: %v\n%s", r, rec, buf)
			err := handlerPanickedError{
				Request: r,
				Err:     rec,
			}
			done(err)
			r.RespondError(err)
			return
		}

		if !responded {
			err := handlerTerminatedError{
				Request: r,
			}
			done(err)
			r.RespondError(err)
		}
	}()

	if err := c.handleRequest(ctx, r, done); err != nil {
		if err == context.Canceled {
			select {
			case <-parentCtx.Done():
				// We canceled the parent context because of an
				// incoming interrupt request, so return EINTR
				// to trigger the right behavior in the client app.
				//
				// Only do this when it's the parent context that was
				// canceled, not a context controlled by the program
				// using this library, so we don't return EINTR too
				// eagerly -- it might cause busy loops.
				//
				// Decent write-up on role of EINTR:
				// http://250bpm.com/blog:12
				err = fuse.EINTR
			default:
				// nothing
			}
		}
		done(err)
		r.RespondError(err)
	}

	// disarm runtime.Goexit protection
	responded = true
}

// handleRequest will either a) call done(s) and r.Respond(s) OR b) return an error.
func (c *Client) handleRequest(ctx context.Context, r fuse.Request, done func(resp interface{})) error {
	cl := c.Client
	Debug("Client:Handle %v", r)
	switch r := r.(type) {
	default:
		// Note: To FUSE, ENOSYS means "this server never implements this request."
		// It would be inappropriate to return ENOSYS for other operations in this
		// switch that might only be unavailable in some contexts, not all.
		return fuse.ENOSYS

	case *fuse.StatfsRequest:
		s := &fuse.StatfsResponse{}
		done(s)
		r.Respond(s)
		return nil

	// Node operations.
	case *fuse.GetattrRequest:
		Debug("Client:Getattr")
		s := &fuse.GetattrResponse{}
		if err := cl.Call("NetFuseServer.Getattr", r, s); err != nil {
			Debug("Client:Getattr err %v", err)
			return err
		}
		done(s)
		r.Respond(s)
		Debug("Client:Getattr done")
		return nil

	case *fuse.SetattrRequest:
		s := &fuse.SetattrResponse{}
		done(s)
		r.Respond(s)
		return nil

	case *fuse.SymlinkRequest:
		s := &fuse.SymlinkResponse{}
		initLookupResponse(&s.LookupResponse)
		done(s)
		r.Respond(s)
		return nil

	case *fuse.ReadlinkRequest:
		return nil

	case *fuse.LinkRequest:
		s := &fuse.LookupResponse{}
		initLookupResponse(s)
		done(s)
		r.Respond(s)
		return nil

	case *fuse.RemoveRequest:
		done(nil)
		r.Respond()
		return nil

	case *fuse.AccessRequest:
		done(nil)
		r.Respond()
		return nil

	case *fuse.LookupRequest:
		var err error
		s := &fuse.LookupResponse{}
		initLookupResponse(s)
		done(s)
		r.Respond(s)
		return err

	case *fuse.MkdirRequest:
		s := &fuse.MkdirResponse{}
		initLookupResponse(&s.LookupResponse)
		done(s)
		r.Respond(s)
		return nil

	case *fuse.OpenRequest:
		s := &fuse.OpenResponse{}
		done(s)
		r.Respond(s)
		return nil

	case *fuse.CreateRequest:
		s := &fuse.CreateResponse{OpenResponse: fuse.OpenResponse{}}
		initLookupResponse(&s.LookupResponse)
		done(s)
		r.Respond(s)
		return nil

	case *fuse.GetxattrRequest:
		s := &fuse.GetxattrResponse{}
		done(s)
		r.Respond(s)
		return nil

	case *fuse.ListxattrRequest:
		s := &fuse.ListxattrResponse{}
		done(s)
		r.Respond(s)
		return nil

	case *fuse.SetxattrRequest:
		done(nil)
		r.Respond()
		return nil

	case *fuse.RemovexattrRequest:
		done(nil)
		r.Respond()
		return nil

	case *fuse.ForgetRequest:
		done(nil)
		r.Respond()
		return nil

	// Handle operations.
	case *fuse.ReadRequest:
		s := &fuse.ReadResponse{Data: make([]byte, 0, r.Size)}
		done(s)
		r.Respond(s)
		return nil

	case *fuse.WriteRequest:
		s := &fuse.WriteResponse{}
		done(s)
		r.Respond(s)
		return fuse.EIO

	case *fuse.FlushRequest:
		done(nil)
		r.Respond()
		return nil

	case *fuse.ReleaseRequest:
		done(nil)
		r.Respond()
		return nil

	case *fuse.DestroyRequest:
		done(nil)
		r.Respond()
		return nil

	case *fuse.RenameRequest:
		done(nil)
		r.Respond()
		return nil

	case *fuse.MknodRequest:
		s := &fuse.LookupResponse{}
		initLookupResponse(s)
		done(s)
		r.Respond(s)
		return nil

	case *fuse.FsyncRequest:
		done(nil)
		r.Respond()
		return nil

	case *fuse.InterruptRequest:
		done(nil)
		r.Respond()
		return nil

		/*	case *FsyncdirRequest:
				return ENOSYS

			case *GetlkRequest, *SetlkRequest, *SetlkwRequest:
				return ENOSYS

			case *BmapRequest:
				return ENOSYS

			case *SetvolnameRequest, *GetxtimesRequest, *ExchangeRequest:
				return ENOSYS
		*/
	}

	panic("not reached")
}

func errstr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

type Client struct {
	C      net.Conn
	M      *fuse.Conn
	Client *rpc.Client

	Debug   func(msg interface{})
	Context func(ctx context.Context, req fuse.Request) context.Context

	// state, protected by meta
	meta sync.Mutex
	req  map[fuse.RequestID]*serveRequest
	node []*serveNode

	wg sync.WaitGroup
}

func iDebug(i interface{}) {
	Debug("Client:%v", i)
}

func New(proto, addr, dir string, options ...fuse.MountOption) (*Client, error) {
	Debug("Client:Dial %s %s", proto, addr)
	c, err := net.Dial(proto, addr)
	if err != nil {
		return nil, err
	}
	// Ping the server. If there's something wrong we don't even
	// want to try the mount.
	var (
		cl  = rpc.NewClient(c)
		arg = &fuse.StatfsRequest{}
		res = &fuse.StatfsResponse{}
	)
	if err := cl.Call("NetFuseServer.Statfs", arg, res); err != nil {
		return nil, fmt.Errorf("New client call to Root:%v", err)
	}
	Debug("Client:Client ping to server: %v", res)
	m, err := fuse.Mount(dir, options...)
	if err != nil {
		return nil, err
	}
	return &Client{C: c, M: m, Client: cl, Debug: iDebug}, nil
}

// Serve serves the FUSE connection by making calls to the server.  It
// returns only when the connection has been closed or an unexpected
// error occurs.
func (c *Client) Serve() error {
	for {
		Debug("Client:serve loop")
		req, err := c.M.ReadRequest()
		Debug("Client:req %v err %v", req, err)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			Debug("Client:Serve %v", req)
			c.serve(req)
		}()
	}
	return nil
}

func (c *Client) Start() error {
	return c.Serve()
}
