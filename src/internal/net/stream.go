// Copyright 2025 Oleg Lodygensky
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions AND
// limitations under the License.
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
