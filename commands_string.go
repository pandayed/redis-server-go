package main

import "strings"

func registerStringCommands() {
	registerCommand("SET", handleSet)
	registerCommand("GET", handleGet)
	registerCommand("INCR", handleIncr)
	registerCommand("DECR", handleDecr)
	registerCommand("EXISTS", handleExists)
	registerCommand("DEL", handleDel)
}

func handleSet(command []string) []byte {
	if err := validateMinArgs(command, 3, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	storeInstance.Set(command[1], command[2])
	return SerializeSimpleString("OK")
}

func handleGet(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	value, exists := storeInstance.Get(command[1])
	if !exists {
		return SerializeNullBulkString()
	}
	return SerializeBulkString(value)
}

func handleIncr(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	num, err := storeInstance.Incr(command[1])
	if err != nil {
		return SerializeError("ERR value is not an integer or out of range")
	}
	return SerializeInteger(num)
}

func handleDecr(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	num, err := storeInstance.Decr(command[1])
	if err != nil {
		return SerializeError("ERR value is not an integer or out of range")
	}
	return SerializeInteger(num)
}

func handleExists(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	exists := storeInstance.Exists(command[1])
	return SerializeInteger(boolToInt(exists))
}

func handleDel(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	deleted := storeInstance.Delete(command[1])
	return SerializeInteger(boolToInt(deleted))
}

