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
