package main

import "strings"

func registerHashCommands() {
	registerCommand("HSET", handleHSet)
	registerCommand("HGET", handleHGet)
	registerCommand("HGETALL", handleHGetAll)
	registerCommand("HDEL", handleHDel)
	registerCommand("HEXISTS", handleHExists)
	registerCommand("HLEN", handleHLen)
}

func handleHSet(command []string) []byte {
	if err := validateMinArgs(command, 4, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	added := storeInstance.HSet(command[1], command[2], command[3])
	return SerializeInteger(added)
}

func handleHGet(command []string) []byte {
	if err := validateMinArgs(command, 3, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	value, exists := storeInstance.HGet(command[1], command[2])
	if !exists {
		return SerializeNullBulkString()
	}
	return SerializeBulkString(value)
}

func handleHGetAll(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	hash := storeInstance.HGetAll(command[1])
	elements := make([][]byte, 0, len(hash)*2)
	for field, value := range hash {
		elements = append(elements, SerializeBulkString(field))
		elements = append(elements, SerializeBulkString(value))
	}
	return SerializeArray(elements)
}

func handleHDel(command []string) []byte {
	if err := validateMinArgs(command, 3, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	deleted := storeInstance.HDel(command[1], command[2:]...)
	return SerializeInteger(deleted)
}

func handleHExists(command []string) []byte {
	if err := validateMinArgs(command, 3, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	exists := storeInstance.HExists(command[1], command[2])
	return SerializeInteger(boolToInt(exists))
}

func handleHLen(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	length := storeInstance.HLen(command[1])
	return SerializeInteger(length)
}

