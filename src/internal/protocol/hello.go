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
package protocol

import (
	"encoding/base64"
	"pqchat/src/internal/pqc"
)

func BuildHello(id *pqc.Identity) (*HelloMessage, []byte, error) {
	msg := &HelloMessage{
		Type:   "HELLO",
		Pseudo: id.Pseudo,
		UserID: id.UserID,
		Pub:    base64.StdEncoding.EncodeToString(id.Pub),
	}

	raw, err := Marshal(msg)
	if err != nil {
		return nil, nil, err
	}

	sig, err := id.Sign(raw)
	if err != nil {
		return nil, nil, err
	}

	msg.Sig = base64.StdEncoding.EncodeToString(sig)
	final, _ := Marshal(msg)

	return msg, final, nil
}
