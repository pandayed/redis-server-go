package commands

import (
	"strings"

	"github.com/pandayed/redis-server-go/resp"
)

func (e *Executor) registerStringCommands() {
	e.registerCommand("SET", e.handleSet)
	e.registerCommand("GET", e.handleGet)
	e.registerCommand("INCR", e.handleIncr)
	e.registerCommand("DECR", e.handleDecr)
	e.registerCommand("EXISTS", e.handleExists)
	e.registerCommand("DEL", e.handleDel)
}

func (e *Executor) handleSet(command []string) []byte {
	if err := validateMinArgs(command, 3, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	e.dataStore.Set(command[1], command[2])
	return resp.SerializeSimpleString("OK")
}

func (e *Executor) handleGet(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	value, exists := e.dataStore.Get(command[1])
	if !exists {
		return resp.SerializeNullBulkString()
	}
	return resp.SerializeBulkString(value)
}

func (e *Executor) handleIncr(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	num, err := e.dataStore.Incr(command[1])
	if err != nil {
		return resp.SerializeError("ERR value is not an integer or out of range")
	}
	return resp.SerializeInteger(num)
}

func (e *Executor) handleDecr(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	num, err := e.dataStore.Decr(command[1])
	if err != nil {
		return resp.SerializeError("ERR value is not an integer or out of range")
	}
	return resp.SerializeInteger(num)
}

func (e *Executor) handleExists(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	exists := e.dataStore.Exists(command[1])
	return resp.SerializeInteger(boolToInt(exists))
}

func (e *Executor) handleDel(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	deleted := e.dataStore.Delete(command[1])
	return resp.SerializeInteger(boolToInt(deleted))
}
