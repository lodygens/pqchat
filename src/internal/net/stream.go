package net

import (
	"encoding/binary"
	"io"

	"github.com/libp2p/go-libp2p/core/network"
)

func WriteFrame(s network.Stream, data []byte) error {
	if err := binary.Write(s, binary.BigEndian, uint16(len(data))); err != nil {
		return err
	}
	_, err := s.Write(data)
	return err
}

func ReadFrame(rd io.Reader) ([]byte, error) {
	var sz uint16
	if err := binary.Read(rd, binary.BigEndian, &sz); err != nil {
		return nil, err
	}
	buf := make([]byte, sz)
	_, err := io.ReadFull(rd, buf)
	return buf, err
}
