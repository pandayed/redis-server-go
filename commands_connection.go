package main

import "strings"

func registerConnectionCommands() {
	registerCommand("PING", handlePing)
	registerCommand("ECHO", handleEcho)
}

func handlePing(command []string) []byte {
	if len(command) == 1 {
		return SerializeSimpleString("PONG")
	}
	return SerializeBulkString(command[1])
}

func handleEcho(command []string) []byte {
	if err := validateMinArgs(command, 2, strings.ToLower(command[0])); err != nil {
		return SerializeError("ERR " + err.Error())
	}
	return SerializeBulkString(command[1])
}

