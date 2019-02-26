// Copyright 2019 the u-root Authors. All rights _
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bdb

import (
//	"text/template"
)

// Constants and header defintions for ChromeOS kernel partition
// boot descriptor block (BDB)

// Size of SHA256 digest in bytes
const SHA256DigestSize = 32

// Size of RSA4096 key data in bytes
const RSA4096KeyDataSize = 1032

// Size of RSA4096 signature in bytes
const RSA4096SigSize = 512

// Size of ECDSA521 key data in bytes = ceil(521/8) * 2
const ECDSA521KeyDataSize = 132

// Size of ECDSA521 signature in bytes = ceil(521/8) * 2
const ECDSA521SigSize = 132

// Size of RSA3072B key data in bytes
const RSA3072BKeyDataSize = 776

// Size of RSA3072B signature in bytes
const RSA3072BSigSize = 384

// Header for

// Magic number for header.magic
const HeaderMagic = 0x30426442

// Current version of header struct
const HeaderVersionMajor = 1
const HeaderVersionMinor = 0

// Expected size of header struct in bytes
const HeaderExpectedSize = 32

type Version struct {
	Major uint8
	Minor uint8
}

type Header struct {
	// Magic number to identify struct = HeaderMagic.
	Magic uint32

	Version

	// Size of structure in bytes
	HeaderSize uint16

	/* Recommended address in SP SRAM to load .  Set to -1 to use
	 * default address. */
	Address uint64

	// Size of the entire  in bytes
	Size uint32

	/* Number of bytes following the  key which are signed by the
	 * header signature. */
	SignedSize uint32

	// Size of OEM area 0 in bytes, or 0 if not present
	OEMArea0Size uint32

	// Reserved; set 0
	_ [8]uint8
}

//***************************************************************************
// Public key structure for

// Magic number for key.magic
const KeyMagic = 0x73334256

// Current version of key struct
const KeyVersionMajor = 1
const KeyVersionMinor = 0

// Supported hash algorithms
const (
	HashAlgInvalid = iota // Not used; invalid
	HashAlgSHA256  = 2
)

// Supported signature algorithms
const (
	SigAlgInvalid   = iota // Not used; invalid
	SigAlgRSA4096   = 3    // RSA-4096, exponent 65537
	SigAlgECSDSA521 = 5    // ECDSA-521
	SigAlgRSA3072B  = 7    // RSA_3072, exponent 3
)

/*
 * Expected size of key struct in bytes, not counting variable-length key
 * data at end.
 */
const KeyEXPECTEDSize = 80

type Key struct {
	// Magic number to identify struct = KeyMagic.
	Magic uint32

	Version

	// Size of structure in bytes, including variable-length key data
	Size uint16

	HashAlg uint8

	SigAlg uint8

	// Reserved; set 0
	_ [2]uint8

	// Key version
	KeyVersion uint32

	// Description; null-terminated ASCII
	Description [128]byte

	/*
	 * Key data.  Variable-length; size is size -
	 * offset_of(key, key_data).
	 */
	Data []uint8
}

//***************************************************************************
// Signature structure for

// Magic number for sig.magic
const SigMagic = 0x6b334256

// Current version of sig struct
const SigVersionMajor = 1
const SigVersionMinor = 0

type Sig struct {
	// Magic number to identify struct = SigMagic.
	Magic uint32

	Version

	/* Size of structure in bytes, including variable-length signature
	 * data. */
	Size uint16

	HashAlg uint8

	SigAlg uint8

	// Reserved; set 0
	_ [2]uint8

	// Number of bytes of data signed by this signature
	SignedSize uint32

	// Description; null-terminated ASCII
	description [128]byte

	/* Signature data.  Variable-length; size is size -
	 * offset_of(sig, sig_data). */
	SigData []uint8
}

//***************************************************************************
// Data structure for

// Magic number for data.magic
const DataMagic = 0x31426442

// Current version of sig struct
const DataVersionMajor = 1
const DataVersionMinor = 0

type Data struct {
	// Magic number to identify struct = DataMagic.
	Magic uint32

	Version

	// Size of structure in bytes, NOT including hashes which follow.
	Size uint16

	// Version of data (RW firmware) contained
	DataVersion uint32

	// Size of OEM area 1 in bytes, or 0 if not present
	OEMArea1Size uint32

	// Number of hashes which follow
	NumHashes uint8

	// Size of each hash entry in bytes
	HashEntrySize uint8

	// Reserved; set 0
	_ [2]uint8

	/* Number of bytes of data signed by the subkey, including this
	 * header */
	SignedSize uint32

	// Reserved; set 0
	_ [8]uint8

	// Description; null-terminated ASCII
	Description [128]byte
}

// Type of data for hash.type
const (
	// Types of data for boot descriptor blocks
	DataSPRW = 1 // SP-RW firmware
	DataAPRW = 2 // AP-RW firmware
	DataMCU  = 3 // MCU firmware

	// Types of data for kernel descriptor blocks
	Kernel   = 128 // Kernel
	CmdLine  = 129 // Command line
	Header16 = 130 // 16-bit vmlinuz header
)

type PartNo uint8
type PartType uint8

// Hash entries which follow the structure
type Hash struct {
	// Offset of data from start of partition
	Offset uint64

	// Size of data in bytes
	Size uint32

	// Partition number containing data
	Partition PartNo

	PartType PartType

	// Reserved; set 0
	_ [2]uint8

	// Address in RAM to load data.  -1 means use default.
	Address uint64

	// SHA-256 hash digest
	Digest [SHA256DigestSize]uint8
}

func (*Header) String() string {
	return ""
}
