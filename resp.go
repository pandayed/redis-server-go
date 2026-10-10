package main

import (
	"bufio"

	"github.com/pandayed/redis-server-go/resp"
)

// These aliases and wrappers preserve the API used by existing tests.
type RESPType = resp.RESPType
type RESPValue = resp.RESPValue

const (
	SimpleString = resp.SimpleString
	Error        = resp.Error
	Integer      = resp.Integer
	BulkString   = resp.BulkString
	Array        = resp.Array
)

func ReadRESP(reader *bufio.Reader) (RESPValue, error) {
	return resp.ReadRESP(reader)
}

func SerializeSimpleString(s string) []byte {
	return resp.SerializeSimpleString(s)
}

func SerializeError(msg string) []byte {
	return resp.SerializeError(msg)
}

func SerializeInteger(n int) []byte {
	return resp.SerializeInteger(n)
}

func SerializeBulkString(s string) []byte {
	return resp.SerializeBulkString(s)
}

func SerializeNullBulkString() []byte {
	return resp.SerializeNullBulkString()
}

func SerializeArray(elements [][]byte) []byte {
	return resp.SerializeArray(elements)
}

func SerializeNullArray() []byte {
	return resp.SerializeNullArray()
}
