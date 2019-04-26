// Hellofs implements a simple "hello world" file system.
package serve

import (
	"testing"
)

func TestNewServer(t *testing.T) {
	// Just create one at 127.0.0.1
	s, err := New("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("New: got %v, want nil", err)
	}
	t.Logf("New server is %v", s)
}
