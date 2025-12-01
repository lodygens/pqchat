package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"pqchat/src/internal/net"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	h, err := net.NewRelayHost()
	if err != nil {
		panic(err)
	}

	fmt.Println("Relay PeerID:", h.ID())
	for _, a := range h.Addrs() {
		fmt.Println(" ", a)
	}

	<-ctx.Done()
	fmt.Println("Shutting down relay…")
	h.Close()
}
