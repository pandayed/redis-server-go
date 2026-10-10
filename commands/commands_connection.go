package commands

import (
	"strings"

	"github.com/pandayed/redis-server-go/resp"
)

func (e *Executor) registerConnectionCommands() {
	e.registerCommand("PING", e.handlePing)
	e.registerCommand("ECHO", e.handleEcho)
}

func (e *Executor) handlePing(command []string) []byte {
	if len(command) == 1 {
		return resp.SerializeSimpleString("PONG")
	}
	return resp.SerializeBulkString(command[1])
}

func (e *Executor) handleEcho(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	return resp.SerializeBulkString(command[1])
}
