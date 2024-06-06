// Copyright 2022 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pci

import (
	"bytes"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPrint(t *testing.T) {
	for _, tt := range []struct {
		name     string
		devices  Devices
		verbose  int
		confSize int
		want     string
	}{
		{
			name: "test1",
			devices: Devices{
				&PCI{
					Config: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
						0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
						0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
						0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
						0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
						0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
						0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
						0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					},
					Bridge: true,
					IO: BAR{
						Base: 64,
					},
					IRQPin: 10,
				},
			},
			verbose:  10,
			confSize: 10,
			want: `: :  
	Control: I/O- Memory- DMA- Special- MemWINV- VGASnoop- ParErr- Stepping- SERR- FastB2B- DisInt-
	Status: INTx- Cap- 66MHz- UDF- FastB2b- ParErr- DEVSEL- DEVSEL=fast <MABORT- >SERR- <PERR-
	Latency: 0, Cache Line Size: 255 bytes
	Bus: primary=00, secondary=00, subordinate=00, sec-latency=
	I/O behind bridge: 0x00000040-0x00000000 [size=0xffffffffffffffc1]
	Memory behind bridge:  [disabled]
	Prefetchable memory behind bridge:  [disabled]
	Interrupt: pin 13 routed to IRQ 0
00000000  ff ff ff ff ff ff ff ff  ff ff 
`,
		},
		{
			name: "test2",
			devices: Devices{
				&PCI{
					Config: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					},
					BARS: []BAR{
						{
							Base: 64,
						},
					},
				},
			},
			verbose:  0,
			confSize: 0,
			want: `: :  
	Control: I/O- Memory- DMA- Special- MemWINV- VGASnoop- ParErr- Stepping- SERR- FastB2B- DisInt-
	Status: INTx- Cap- 66MHz- UDF- FastB2b- ParErr- DEVSEL- DEVSEL=fast <MABORT- >SERR- <PERR-
	Latency: 0
	Region 0: Memory at 00000040 (32-bit, non-prefetchable) [size=0xffffffffffffffc1]

`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var writeBuf = &bytes.Buffer{}
			if got := tt.devices.Print(writeBuf, 10, tt.confSize); got == nil {
				if writeBuf.String() != tt.want {
					t.Errorf("Buffer contains = %q, want: %q", writeBuf.String(), tt.want)
				}
			}
		})
	}
}

func TestSetVendorDeviceName(t *testing.T) {
	for _, tt := range []struct {
		name           string
		devices        Devices
		VendorNameWant string
		DeviceNameWant string
	}{
		{
			name: "Lookup Using ID 80ee Device cafe",
			devices: Devices{
				&PCI{
					Vendor: 0x80ee,
					Device: 0xcafe,
				},
			},
			VendorNameWant: "InnoTek Systemberatung GmbH",
			DeviceNameWant: "VirtualBox Guest Service",
		},
		{
			name: "Lookup Using ID 1055 Device e420",
			devices: Devices{
				&PCI{
					Vendor: 0x1055,
					Device: 0xe420,
				},
			},
			VendorNameWant: "Efar Microsystems",
			DeviceNameWant: "LAN9420/LAN9420i",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tt.devices.SetVendorDeviceName()
			VendorNameGot, DeviceNameGot := tt.devices[0].VendorName, tt.devices[0].DeviceName
			if VendorNameGot != tt.VendorNameWant {
				t.Errorf("Vendor mismatch, got: %q, want: %q\n", VendorNameGot, tt.VendorNameWant)
			}
			if DeviceNameGot != tt.DeviceNameWant {
				t.Errorf("Device mismatch, got: %q, want: %q\n", DeviceNameGot, tt.DeviceNameWant)
			}
		})
	}
}

func TestReadConfig(t *testing.T) {
	configBytes := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}
	dir := t.TempDir()
	f, err := os.Create(filepath.Join(dir, "config"))
	if err != nil {
		t.Errorf("Creating file failed: %v", err)
	}
	_, err = f.Write(configBytes)
	if err != nil {
		t.Errorf("Writing to file failed: %v", err)
	}
	for _, tt := range []struct {
		name        string
		devices     Devices
		controlWant Control
		statusWant  Status
		errWant     string
	}{
		{
			name: "Reading config file",
			devices: Devices{
				&PCI{
					FullPath: dir,
				},
			},
			controlWant: 0x5544,
			statusWant:  0x7766,
		},
		{
			name: "config file does not exist",
			devices: Devices{
				&PCI{
					FullPath: "d",
				},
			},
			errWant: "no such file or directory",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.devices.ReadConfig(); got != nil {
				if !strings.Contains(got.Error(), tt.errWant) {
					t.Errorf("ReadConfig() = %q, want to contain: %q", got, tt.errWant)
				}
			} else {
				if tt.devices[0].Control != tt.controlWant {
					t.Errorf("Control is = '%#x', want: '%#x'", tt.devices[0].Control, tt.controlWant)
				}
				if tt.devices[0].Status != tt.statusWant {
					t.Errorf("Status is = '%#x', want: '%#x'", tt.devices[0].Status, tt.statusWant)
				}
			}
		})
	}
}

func TestReadConfigRegister(t *testing.T) {
	configBytes := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}
	dir := t.TempDir()
	f, err := os.Create(filepath.Join(dir, "config"))
	if err != nil {
		t.Errorf("Creating file failed: %v", err)
	}
	_, err = f.Write(configBytes)
	if err != nil {
		t.Errorf("Writing to file failed: %v", err)
	}
	for _, tt := range []struct {
		name     string
		devices  Devices
		offset   int64
		size     int64
		valsWant uint64
		errWant  string
	}{
		{
			name: "read byte 2 from config file",
			devices: Devices{
				&PCI{
					FullPath: dir,
				},
			},
			offset:   1,
			size:     8,
			valsWant: 0x11,
		},
		{
			name: "read byte 2 & 3 from config file",
			devices: Devices{
				&PCI{
					FullPath: dir,
				},
			},
			offset:   1,
			size:     16,
			valsWant: 0x2211,
		},
		{
			name: "wrong size",
			devices: Devices{
				&PCI{
					FullPath: dir,
				},
			},
			offset:  0,
			size:    0,
			errWant: "only options are 8, 16, 32, 64",
		},
		{
			name: "config file does not exist",
			devices: Devices{
				&PCI{
					FullPath: "d",
				},
			},
			errWant: "no such file or directory",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if vals, got := tt.devices.ReadConfigRegister(tt.offset, tt.size); got != nil {
				if !strings.Contains(got.Error(), tt.errWant) {
					t.Errorf("ReadConfig() = %q, want to contain: %q", got, tt.errWant)
				}
			} else {
				if vals[0] != tt.valsWant {
					t.Errorf("ReadConfig() = '%#x', want: '%#x'", vals[0], tt.valsWant)
				}
			}
		})
	}
}

func TestWriteConfigRegister(t *testing.T) {
	configBytes := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77}
	dir := t.TempDir()
	f, err := os.Create(filepath.Join(dir, "config"))
	if err != nil {
		t.Errorf("Creating file failed: %v", err)
	}
	_, err = f.Write(configBytes)
	if err != nil {
		t.Errorf("Writing to file failed: %v", err)
	}
	for _, tt := range []struct {
		name    string
		devices Devices
		offset  int64
		size    int64
		val     uint64
		want    string
		errWant string
	}{
		{
			name: "Writing 1 byte to config file with offset 1",
			devices: Devices{
				&PCI{
					FullPath: dir,
				},
			},
			offset: 1,
			size:   8,
			val:    0x22,
			want:   "0022223344556677",
		},
		{
			name: "Writing 2 bytes to config file with offset 1",
			devices: Devices{
				&PCI{
					FullPath: dir,
				},
			},
			offset: 1,
			size:   16,
			val:    0x4433,
			want:   "0033443344556677",
		},
		{
			name: "config file does not exist",
			devices: Devices{
				&PCI{
					FullPath: "d",
				},
			},
			errWant: "no such file or directory",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.devices.WriteConfigRegister(tt.offset, tt.size, tt.val); got != nil {
				if !strings.Contains(got.Error(), tt.errWant) {
					t.Errorf("ReadConfig() = %q, want to contain: %q", got, tt.errWant)
				}
			} else {
				got, err := os.ReadFile(filepath.Join(dir, "config"))
				if err != nil {
					t.Errorf("Failed to read file %v", err)
				}
				if hex.EncodeToString(got) != tt.want {
					t.Errorf("Config file contains = %q, want: %q", hex.EncodeToString(got), tt.want)
				}
			}
		})
	}
}
