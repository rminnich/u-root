// Hellofs implements a simple "hello world" file system.
package serve

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"bazil.org/fuse"
	"github.com/u-root/u-root/pkg/netfuse/serve"
	"golang.org/x/sys/unix"
)

func TestNewClient(t *testing.T) {
	Debug = t.Logf
	fs, err := serve.NewFileSystem("/tmp")
	if err != nil {
		t.Fatalf("NewFileSystem: got %v, want nil", err)
	}
	// Just create one at 127.0.0.1
	s, err := serve.New("tcp", "127.0.0.1:0", fs)
	if err != nil {
		t.Fatalf("New: got %v, want nil", err)
	}
	t.Logf("New server is %v", s)
	d, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		if err := s.Start(); err != nil {
			t.Fatalf("Start server: got %v,want nil", err)
		}
		if err := s.Run(); err != nil {
			t.Fatalf("Run server: got %v, want nil", err)
		}
	}()
	// Talk to me.
	a := s.Server.Addr()
	t.Logf("Server is listening on %v", a)
	cl, err := New(a.Network(), a.String(), d)
	if err != nil {
		t.Fatalf("New Client: got %v, want nil", err)
	}
	defer fuse.Unmount(d)
	go func() {
		if err := cl.Start(); err != nil {
			t.Fatalf("Start Client: got %v, want nil", err)
		}
	}()
	t.Logf("Now stat things")
	fi, err := os.Stat(d)
	if err != nil {
		t.Fatalf("Stat mount point %q: %v", d, err)
	}
	t.Logf("Stat %q returns %v", d, fi)

}

func TestMain(m *testing.M) {
	if false { // still doesn't fix it. Run me under unshare.
		if err := unix.Unshare(unix.CLONE_NEWNS); err != nil {
			log.Fatal(err)
		}
	}

	retCode := m.Run()
	os.Exit(retCode)
}
