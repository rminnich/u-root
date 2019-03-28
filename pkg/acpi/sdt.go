// Copyright 2019 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package acpi

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type SDT struct {
	Generic
	Tables []int64
	Base   int64
}

func init() {
	addUnMarshaler("RSDT", unmarshalSDT)
	addUnMarshaler("XSDT", unmarshalSDT)
}

func unmarshalSDT(t Tabler) (Tabler, error) {
	s := &SDT{
		Generic: Generic{
			Header: *GetHeader(t),
			data:   t.AllData(),
		},
	}

	sig := s.Sig()
	if sig != "RSDT" && sig != "XSDT" {
		return nil, fmt.Errorf("%v is not RSDT or XSDT", sig)
	}

	// Now the fun. In 1999, 64-bit micros had been out for about 10 years.
	// Intel had announced the ia64 years earlier. In 2000 the ACPI committee
	// chose 32-bit pointers anyway, then had to backfill a bunch of table
	// types to do 64 bits. Geez.
	esize := 4
	if sig == "XSDT" {
		esize = 8
	}
	d := t.TableData()

	for i := 0; i < len(d); i += esize {
		val := int64(0)
		if sig == "XSDT" {
			val = int64(binary.LittleEndian.Uint64(d[i : i+8]))
		} else {
			val = int64(binary.LittleEndian.Uint32(d[i : i+4]))
		}
		s.Tables = append(s.Tables, val)
	}
	return s, nil
}

func (s *SDT) Marshal() ([]byte, error) {
	h, err := s.Generic.Header.Marshal()
	if err != nil {
		return nil, err
	}
	b := bytes.NewBuffer(h)
	x := s.Sig() == "XSDT"
	for _, p := range s.Tables {
		if x {
			binary.Write(b, binary.LittleEndian, p)
		} else {
			binary.Write(b, binary.LittleEndian, uint32(p))
		}
	}
	return b.Bytes(), nil
}

// MarshalAll marshals out an SDT, and all the tables, in a blob suitable for
// kexec. All addresses are recomputed, as there may be more tables. Further,
// even if tables were scattered all over, we unify them into one segment.
// We're going to try an experiment here: no matter whether we pulled in an RSDT
// or XSDT, we're going to always write an XSDT. Easy to change later and what can
// go wrong? Besides the ACPICA code shitting the bed, of course, which it does
// have a habit of doing.
// The basic idea:
// 1. Count the number of tables
// 2. compute the new XSDT size.
// 3. Serialize the tables.
// 4. Compute []int64 addresses
// 5. Serialize out the SDT header
// 6. addresses
// 7. tables
func (s *SDT) MarshalAll(t ...Tabler) ([]byte, error) {
	var tabs [][]byte
	for _, addr := range s.Tables {
		// We need to read the table in, and add it the things
		// we marshal out. That's far safer than assuming it just ends up magically
		// there.
		t, err := ReadRaw(addr)
		if err != nil {
			return nil, err
		}
		tabs = append(tabs, t.AllData())
	}

	for _, tt := range t {
		b, err := Marshal(tt)
		if err != nil {
			return nil, err
		}
		tabs = append(tabs, b)
	}

	Debug("processed tables")
	// The length of the SDT is SSDTSize + len(s.Tables) * 8 (64-bit pointers)
	// The easiest path here is to replace the data with the new data, but first we have to
	// compute the pointers. So we do this as follows:
	// truncate ssd to just the header.
	s.Generic.data = s.Generic.data[:SSDTSize]
	var (
		addrs bytes.Buffer
		st    []byte
	)

	base := s.Base
	for _, t := range tabs {
		st = append(st, t...)
		binary.Write(&addrs, binary.LittleEndian, base)
		base += int64(len(t))
	}
	s.Generic.data = append(s.Generic.data, append(addrs.Bytes(), st...)...)
	h, err := s.Generic.Marshal()
	if err != nil {
		return nil, err
	}

	return h, nil
}

func ReadSDT() (*SDT, error) {
	r, err := GetRSDP()
	if err != nil {
		return nil, err
	}
	s, err := UnMarshalSDT(r)
	return s, err
}

func NewSDT(opt ...func(*SDT)) *SDT {
	var s = &SDT{
		Generic: Generic{
			Header: Header{
				Sig:             "XSDT",
				Length:          SSDTSize,
				Revision:        1,
				OEMID:           "GOOGLE",
				OEMTableID:      "ACPI=TOY",
				OEMRevision:     1,
				CreatorID:       1,
				CreatorRevision: 1,
			},
			data: make([]byte, SSDTSize),
		},
	}
	for _, o := range opt {
		o(s)
	}
	return s
}
