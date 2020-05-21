package qcow2

import (
	"encoding/binary"
	"io"
)

// ported by cloudius-systems/capstan
// https://github.com/cloudius-systems/capstan/blob/bf6c68a40f6d70a1c7bb457dbf1bc0b87b69d2b3/image/qcow2/qcow2.go

const (
	// Qcow2Magic is magic header of qcow2 file
	Qcow2Magic = ('Q' << 24) | ('F' << 16) | ('I' << 8) | 0xfb
)

// Header is qcow2 file header
type Header struct {
	Magic                 uint32
	Version               uint32
	BackingFileOffset     uint64
	BackingFileSize       uint32
	ClusterBits           uint32
	Size                  uint64
	CryptMethod           uint32
	L1Size                uint32
	L1TableOffset         uint64
	RefcountTableOffset   uint64
	RefcountTableClusters uint32
	NbSnapshots           uint32
	SnapshotsOffset       uint64
}

// Probe is validator of qcow2
func Probe(r io.Reader) (bool, *Header) {
	header, err := readHeader(r)
	if err != nil {
		return false, nil
	}
	return header.Magic == Qcow2Magic, header
}

func readHeader(r io.Reader) (*Header, error) {
	var header Header
	err := binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return nil, err
	}
	return &header, nil
}
