#!/usr/bin/env bash

set -e

which tmux > /dev/null || { echo "tmux not found"; exit 1; }

# Colors
GREEN="\033[1;32m"
BLUE="\033[1;34m"
RED="\033[1;31m"
NC="\033[0m"

echo -e "${GREEN}=== PQChat Local Test Script ===${NC}"

# Paths
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_DIR="$ROOT_DIR/bin"

PQCHAT="$BIN_DIR/pqchat"
RELAYER="$BIN_DIR/relayer"

if [[ ! -x "$PQCHAT" ]]; then
    echo -e "${RED}ERROR:${NC} $PQCHAT not found or not executable"
    exit 1
fi

if [[ ! -x "$RELAYER" ]]; then
    echo -e "${RED}ERROR:${NC} $RELAYER not found or not executable"
    exit 1
fi

SESSION="pqchat-test"

# Kill old session if exists
tmux has-session -t $SESSION 2>/dev/null && tmux kill-session -t $SESSION

echo -e "${BLUE}Launching TMUX session '$SESSION'...${NC}"

tmux new-session -d -s $SESSION -n relay

# 1) Relay
tmux send-keys -t $SESSION:relay "clear && echo '[RELAY]' && $RELAYER" C-m

# Split for Alice
tmux split-window -v -t $SESSION:relay
tmux select-pane -t 1
tmux send-keys "clear && echo '[ALICE]' && $PQCHAT -pseudo alice -relay /ip4/127.0.0.1/tcp/4001/p2p/\$(grep -m1 'p2p/' <(sleep 1 && $RELAYER --print-id))" C-m

# Split for Bob
tmux split-window -h -t $SESSION:relay
tmux select-pane -t 2
tmux send-keys "clear && echo '[BOB]' && $PQCHAT -pseudo bob -relay /ip4/127.0.0.1/tcp/4001/p2p/\$(grep -m1 'p2p/' <(sleep 1 && $RELAYER --print-id))" C-m

# Split for Charlie
tmux split-window -v -t $SESSION:relay
tmux select-pane -t 3
tmux send-keys "clear && echo '[CHARLIE]' && $PQCHAT -pseudo charlie -relay /ip4/127.0.0.1/tcp/4001/p2p/\$(grep -m1 'p2p/' <(sleep 1 && $RELAYER --print-id))" C-m

sleep 2

echo -e "${BLUE}Retrieving Alice's multiaddr...${NC}"
ALICE_ADDR=$(tmux capture-pane -pt $SESSION:relay.1 | grep -m1 "/ip4" | awk '{print $1}')
echo -e "${GREEN}Alice addr: $ALICE_ADDR${NC}"

sleep 1

echo -e "${BLUE}Connecting Bob and Charlie to Alice...${NC}"
tmux send-keys -t $SESSION:relay.2 "/connect $ALICE_ADDR" C-m
tmux send-keys -t $SESSION:relay.3 "/connect $ALICE_ADDR" C-m

sleep 2

echo -e "${BLUE}Sending test messages...${NC}"

# Alice says hello
tmux send-keys -t $SESSION:relay.1 "Hello everyone! I am Alice." C-m

# Bob whispers to Alice
tmux send-keys -t $SESSION:relay.2 "@alice-secret Hello Alice!" C-m

# Charlie sends a broadcast
tmux send-keys -t $SESSION:relay.3 "Hi all — Charlie here." C-m

echo -e "${GREEN}=== PQChat local test launched! ===${NC}"
echo -e "${GREEN}Type 'tmux attach -t pqchat-test' to view.${NC}"
