package commands

import (
	"strings"

	"github.com/pandayed/redis-server-go/resp"
)

func (e *Executor) registerHashCommands() {
	e.registerCommand("HSET", e.handleHSet)
	e.registerCommand("HGET", e.handleHGet)
	e.registerCommand("HGETALL", e.handleHGetAll)
	e.registerCommand("HDEL", e.handleHDel)
	e.registerCommand("HEXISTS", e.handleHExists)
	e.registerCommand("HLEN", e.handleHLen)
}

func (e *Executor) handleHSet(command []string) []byte {
	if err := validateMinArgs(command, 4, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	added := e.dataStore.HSet(command[1], command[2], command[3])
	return resp.SerializeInteger(added)
}

func (e *Executor) handleHGet(command []string) []byte {
	if err := validateMinArgs(command, 3, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	value, exists := e.dataStore.HGet(command[1], command[2])
	if !exists {
		return resp.SerializeNullBulkString()
	}
	return resp.SerializeBulkString(value)
}

func (e *Executor) handleHGetAll(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	hash := e.dataStore.HGetAll(command[1])
	elements := make([][]byte, 0, len(hash)*2)
	for field, value := range hash {
		elements = append(elements, resp.SerializeBulkString(field))
		elements = append(elements, resp.SerializeBulkString(value))
	}
	return resp.SerializeArray(elements)
}

func (e *Executor) handleHDel(command []string) []byte {
	if err := validateMinArgs(command, 3, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	deleted := e.dataStore.HDel(command[1], command[2:]...)
	return resp.SerializeInteger(deleted)
}

func (e *Executor) handleHExists(command []string) []byte {
	if err := validateMinArgs(command, 3, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	exists := e.dataStore.HExists(command[1], command[2])
	return resp.SerializeInteger(boolToInt(exists))
}

func (e *Executor) handleHLen(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	length := e.dataStore.HLen(command[1])
	return resp.SerializeInteger(length)
}
