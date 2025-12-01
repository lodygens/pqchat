package chat

import (
	"encoding/json"
	"fmt"
	"pqchat/src/internal/protocol"
)

func HandleIncoming(raw []byte) {
	// Try Hello
	var hello protocol.HelloMessage
	if json.Unmarshal(raw, &hello) == nil && hello.Type == "HELLO" {
		fmt.Println("[+] HELLO from", hello.Pseudo)

		// TODO: verify signature ML-DSA

		return
	}

	// Try Chat
	var chat protocol.ChatMessage
	if json.Unmarshal(raw, &chat) == nil && chat.Type == "CHAT" {
		HandleChat(&chat)
		return
	}

	fmt.Println("[?] Unknown message")
}
