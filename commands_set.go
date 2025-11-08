package main

import "strings"

func registerSetCommands() {
	registerCommand("SADD", handleSAdd)
	registerCommand("SMEMBERS", handleSMembers)
	registerCommand("SISMEMBER", handleSIsMember)
	registerCommand("SREM", handleSRem)
	registerCommand("SCARD", handleSCard)
}

func handleSAdd(command []string) []byte {
	if err := validateMinArgs(command, 3, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	added := storeInstance.SAdd(command[1], command[2:]...)
	return SerializeInteger(added)
}

func handleSMembers(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	members := storeInstance.SMembers(command[1])
	return serializeStringArray(members)
}

func handleSIsMember(command []string) []byte {
	if err := validateMinArgs(command, 3, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	isMember := storeInstance.SIsMember(command[1], command[2])
	return SerializeInteger(boolToInt(isMember))
}

func handleSRem(command []string) []byte {
	if err := validateMinArgs(command, 3, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	removed := storeInstance.SRem(command[1], command[2:]...)
	return SerializeInteger(removed)
}

func handleSCard(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	cardinality := storeInstance.SCard(command[1])
	return SerializeInteger(cardinality)
}

