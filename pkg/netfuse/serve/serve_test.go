// Hellofs implements a simple "hello world" file system.
package serve

import (
	"log"
	"net"
	"net/rpc"
	"testing"

	"bazil.org/fuse"
)

func TestNewServer(t *testing.T) {
	Debug = t.Logf
	fs, err := NewFileSystem("/tmp")
	if err != nil {
		t.Fatalf("NewFileSystem: got %v, want nil", err)
	}
	// Just create one at 127.0.0.1
	s, err := New("tcp", "127.0.0.1:0", fs)
	if err != nil {
		t.Fatalf("New: got %v, want nil", err)
	}
	t.Logf("New server is %v", s)
	// Talk to me.
	a := s.Server.Addr()
	c, err := net.Dial(a.Network(), a.String())
	if err != nil {
		t.Fatalf("Dial: got %v, want nil", err)
	}
	cl := rpc.NewClient(c)
	log.Printf("client %v on conn %v", cl, c)
	var (
		arg = &fuse.StatfsRequest{}
		res = &fuse.StatfsResponse{}
	)
	if err := s.Start(); err != nil {
		t.Fatalf("Server start: got %v, want nil", err)
	}
	go func() {
		s.Run()
	}()
	if err := cl.Call("NetFuseServer.Statfs", arg, res); err != nil {
		t.Fatalf("statfs: got %v, want nil", err)
	}

}
