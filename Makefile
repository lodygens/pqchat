GO              := go
BIN_DIR         := bin
SRC_DIR         := src

RELAYER_BIN     := $(BIN_DIR)/relayer
CHAT_BIN        := $(BIN_DIR)/pqchat

.PHONY: all
all: relayer chat

.PHONY: relayer
relayer:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(RELAYER_BIN) ./$(SRC_DIR)/cmd/relayer

.PHONY: chat
chat:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(CHAT_BIN) ./$(SRC_DIR)/cmd/pqchat

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)

.PHONY: run-relayer
run-relayer: relayer
	./$(RELAYER_BIN)

.PHONY: run-chat
run-chat: chat
	./$(CHAT_BIN)

.PHONY: local-test
local-test: chat relayer
	scripts/local-test.sh