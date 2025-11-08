package main

import "strings"

func registerListCommands() {
	registerCommand("LPUSH", handleLPush)
	registerCommand("RPUSH", handleRPush)
	registerCommand("LPOP", handleLPop)
	registerCommand("RPOP", handleRPop)
	registerCommand("LRANGE", handleLRange)
	registerCommand("LLEN", handleLLen)
}

func handleLPush(command []string) []byte {
	if err := validateMinArgs(command, 3, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	length := storeInstance.LPush(command[1], command[2:]...)
	return SerializeInteger(length)
}

func handleRPush(command []string) []byte {
	if err := validateMinArgs(command, 3, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	length := storeInstance.RPush(command[1], command[2:]...)
	return SerializeInteger(length)
}

func handleLPop(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	value, exists := storeInstance.LPop(command[1])
	if !exists {
		return SerializeNullBulkString()
	}
	return SerializeBulkString(value)
}

func handleRPop(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	value, exists := storeInstance.RPop(command[1])
	if !exists {
		return SerializeNullBulkString()
	}
	return SerializeBulkString(value)
}

func handleLRange(command []string) []byte {
	if err := validateMinArgs(command, 4, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	
	intArgs, err := parseIntArgs(command[2:4])
	if err != nil {
		return SerializeError("ERR " + err.Error())
	}
	
	values := storeInstance.LRange(command[1], intArgs[0], intArgs[1])
	return serializeStringArray(values)
}

func handleLLen(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	length := storeInstance.LLen(command[1])
	return SerializeInteger(length)
}

func serializeStringArray(values []string) []byte {
	elements := make([][]byte, len(values))
	for i, v := range values {
		elements[i] = SerializeBulkString(v)
	}
	return SerializeArray(elements)
}

