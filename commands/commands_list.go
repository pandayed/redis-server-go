package commands

import (
	"strings"

	"github.com/pandayed/redis-server-go/resp"
)

func (e *Executor) registerListCommands() {
	e.registerCommand("LPUSH", e.handleLPush)
	e.registerCommand("RPUSH", e.handleRPush)
	e.registerCommand("LPOP", e.handleLPop)
	e.registerCommand("RPOP", e.handleRPop)
	e.registerCommand("LRANGE", e.handleLRange)
	e.registerCommand("LLEN", e.handleLLen)
}

func (e *Executor) handleLPush(command []string) []byte {
	if err := validateMinArgs(command, 3, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	length := e.dataStore.LPush(command[1], command[2:]...)
	return resp.SerializeInteger(length)
}

func (e *Executor) handleRPush(command []string) []byte {
	if err := validateMinArgs(command, 3, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	length := e.dataStore.RPush(command[1], command[2:]...)
	return resp.SerializeInteger(length)
}

func (e *Executor) handleLPop(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	value, exists := e.dataStore.LPop(command[1])
	if !exists {
		return resp.SerializeNullBulkString()
	}
	return resp.SerializeBulkString(value)
}

func (e *Executor) handleRPop(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	value, exists := e.dataStore.RPop(command[1])
	if !exists {
		return resp.SerializeNullBulkString()
	}
	return resp.SerializeBulkString(value)
}

func (e *Executor) handleLRange(command []string) []byte {
	if err := validateMinArgs(command, 4, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}

	intArgs, err := parseIntArgs(command[2:4])
	if err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}

	values := e.dataStore.LRange(command[1], intArgs[0], intArgs[1])
	return serializeStringArray(values)
}

func (e *Executor) handleLLen(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	length := e.dataStore.LLen(command[1])
	return resp.SerializeInteger(length)
}

func serializeStringArray(values []string) []byte {
	elements := make([][]byte, len(values))
	for i, v := range values {
		elements[i] = resp.SerializeBulkString(v)
	}
	return resp.SerializeArray(elements)
}
