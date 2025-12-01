package chat

import (
	"fmt"
	"pqchat/src/internal/protocol"
)

func HandleChat(msg *protocol.ChatMessage) {
	fmt.Printf("<%s> %s\n", msg.From, msg.Body)
}
