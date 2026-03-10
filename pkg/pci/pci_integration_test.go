// Copyright 2012-2026 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/u-root/cpu/client"
	"github.com/u-root/cpu/vm"
)

// TestCPUAMD64 tests both general and specific things. The specific parts are the io and cmos commands.
// It being cheaper to use a single generated initramfs, we use the full u-root for several tests.
func TestCPUAMD64(t *testing.T) {
	t.Logf("TestCPUAMD64 here ...")
	d := t.TempDir()
	i, err := vm.New("linux", "amd64")
	if !errors.Is(err, nil) {
		t.Fatalf("Testing kernel=linux arch=amd64: got %v, want nil", err)
	}

	// Cancel before wg.Wait(), so goroutine can exit.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: have a one-time helper that builds a full u-root image once,
	// that all tests can use.
	// TODO: for all the tests, we need start only one VM. Even for kexec,
	// since it just starts a new kernel, and we can have that kernel use
	// the initramfs that runs cpud.
	n, err := i.Uroot(d)
	if err != nil {
		t.Skipf("skipping this test as we have no uroot command:%v", err)
	}

	c, err := i.CommandContext(ctx, d, n)
	if err != nil {
		t.Fatalf("starting VM: got %v, want nil", err)
	}
	if err := i.StartVM(c); err != nil {
		t.Fatalf("starting VM: got %v, want nil", err)
	}
	log.Printf("you have 5 minutes to try some stuff")
	time.Sleep(5 * time.Minute)
	cpu, err := i.CPUCommand("go", "test", "TestPCI")
	if err != nil {
		t.Fatalf("CPUCommand: got %v, want nil", err)
	}
	client.SetVerbose(t.Logf)

	b, err := cpu.CombinedOutput()
	if err != nil {
		t.Errorf("go test -test.Run TestPCI: got %v, want nil", err)
	}
	t.Logf("output %s", string(b))
}
