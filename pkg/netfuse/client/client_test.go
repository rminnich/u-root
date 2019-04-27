// Hellofs implements a simple "hello world" file system.
package serve

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/u-root/u-root/pkg/netfuse/serve"
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
	// Talk to me.
	a := s.Server.Addr()
	cl, err := New(a.Network(), a.String(), d)
	if err != nil {
		t.Fatalf("New Client: got %v, want nil", err)
	}
	go func() {
		if err := cl.Start(); err != nil {
			t.Fatalf("Start Client: got %v, want nil", err)
		}
	}()
	fi, err := os.Stat(d)
	if err != nil {
		t.Fatalf("Stat mount point %q: %v", d, err)
	}
	t.Logf("Stat %q returns %v", d, fi)

}
