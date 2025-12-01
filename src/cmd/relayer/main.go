// Copyright 2025 Oleg Lodygensky
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions AND
// limitations under the License.

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
