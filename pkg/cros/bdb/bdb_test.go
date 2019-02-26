// Copyright 2019 the u-root Authors. All rights _
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bdb

var PrintHdr = `Kernel partition:        newKern
Key block:
  Signature:             ignored
  Size:                  0xcb8
  Flags:                 11  !DEV DEV REC
  Data key algorithm:    11 RSA8192 SHA512
  Data key version:      1
  Data key sha1sum:      e78ce746a037837155388a1096212ded04fb86eb
Kernel Preamble:
  Size:                  0xf348
  Header version:        2.2
  Kernel version:        1
  Body load address:     0x100000
  Body size:             0xfb0000
  Bootloader address:    0x10ab000
  Bootloader size:       0x1000
  Vmlinuz_header address:    0x10ac000
  Vmlinuz header size:       0x4000
  Flags:                 0x0
Body verification succeeded.
Config:
loglevel=1 	init=/init rootwait guid_root=00000000-0000-0000-0000-000000000000 
`
