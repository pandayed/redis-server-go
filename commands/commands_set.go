package commands

import (
	"strings"

	"github.com/pandayed/redis-server-go/resp"
)

func (e *Executor) registerSetCommands() {
	e.registerCommand("SADD", e.handleSAdd)
	e.registerCommand("SMEMBERS", e.handleSMembers)
	e.registerCommand("SISMEMBER", e.handleSIsMember)
	e.registerCommand("SREM", e.handleSRem)
	e.registerCommand("SCARD", e.handleSCard)
}

func (e *Executor) handleSAdd(command []string) []byte {
	if err := validateMinArgs(command, 3, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	added := e.dataStore.SAdd(command[1], command[2:]...)
	return resp.SerializeInteger(added)
}

func (e *Executor) handleSMembers(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	members := e.dataStore.SMembers(command[1])
	return serializeStringArray(members)
}

func (e *Executor) handleSIsMember(command []string) []byte {
	if err := validateMinArgs(command, 3, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	isMember := e.dataStore.SIsMember(command[1], command[2])
	return resp.SerializeInteger(boolToInt(isMember))
}

func (e *Executor) handleSRem(command []string) []byte {
	if err := validateMinArgs(command, 3, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	removed := e.dataStore.SRem(command[1], command[2:]...)
	return resp.SerializeInteger(removed)
}

func (e *Executor) handleSCard(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return resp.SerializeError("ERR " + err.Error())
	}
	cardinality := e.dataStore.SCard(command[1])
	return resp.SerializeInteger(cardinality)
}
