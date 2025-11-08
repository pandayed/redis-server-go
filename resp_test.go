package main

import (
	"bufio"
	"bytes"
	"testing"
)

func TestReadRESP_SimpleString(t *testing.T) {
	input := "+OK\r\n"
	reader := bufio.NewReader(bytes.NewReader([]byte(input)))

	value, err := ReadRESP(reader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if value.Type != SimpleString {
		t.Errorf("Expected SimpleString type, got %c", value.Type)
	}

	if value.Str != "OK" {
		t.Errorf("Expected 'OK', got '%s'", value.Str)
	}
}

func TestReadRESP_Error(t *testing.T) {
	input := "-ERR unknown command\r\n"
	reader := bufio.NewReader(bytes.NewReader([]byte(input)))

	value, err := ReadRESP(reader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if value.Type != Error {
		t.Errorf("Expected Error type, got %c", value.Type)
	}

	if value.Str != "ERR unknown command" {
		t.Errorf("Expected 'ERR unknown command', got '%s'", value.Str)
	}
}

func TestReadRESP_Integer(t *testing.T) {
	input := ":42\r\n"
	reader := bufio.NewReader(bytes.NewReader([]byte(input)))

	value, err := ReadRESP(reader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if value.Type != Integer {
		t.Errorf("Expected Integer type, got %c", value.Type)
	}

	if value.Num != 42 {
		t.Errorf("Expected 42, got %d", value.Num)
	}
}

func TestReadRESP_BulkString(t *testing.T) {
	input := "$5\r\nhello\r\n"
	reader := bufio.NewReader(bytes.NewReader([]byte(input)))

	value, err := ReadRESP(reader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if value.Type != BulkString {
		t.Errorf("Expected BulkString type, got %c", value.Type)
	}

	if value.Bulk != "hello" {
		t.Errorf("Expected 'hello', got '%s'", value.Bulk)
	}
}

func TestReadRESP_BulkString_WithSpaces(t *testing.T) {
	input := "$11\r\nhello world\r\n"
	reader := bufio.NewReader(bytes.NewReader([]byte(input)))

	value, err := ReadRESP(reader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if value.Bulk != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", value.Bulk)
	}
}

func TestReadRESP_NullBulkString(t *testing.T) {
	input := "$-1\r\n"
	reader := bufio.NewReader(bytes.NewReader([]byte(input)))

	value, err := ReadRESP(reader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if value.Type != BulkString {
		t.Errorf("Expected BulkString type, got %c", value.Type)
	}

	if !value.Null {
		t.Errorf("Expected null bulk string")
	}
}

func TestReadRESP_Array(t *testing.T) {
	input := "*2\r\n$3\r\nGET\r\n$5\r\nmykey\r\n"
	reader := bufio.NewReader(bytes.NewReader([]byte(input)))

	value, err := ReadRESP(reader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if value.Type != Array {
		t.Errorf("Expected Array type, got %c", value.Type)
	}

	if len(value.Array) != 2 {
		t.Errorf("Expected array length 2, got %d", len(value.Array))
	}

	if value.Array[0].Bulk != "GET" {
		t.Errorf("Expected 'GET', got '%s'", value.Array[0].Bulk)
	}

	if value.Array[1].Bulk != "mykey" {
		t.Errorf("Expected 'mykey', got '%s'", value.Array[1].Bulk)
	}
}

func TestReadRESP_ComplexArray(t *testing.T) {
	input := "*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$7\r\nmyvalue\r\n"
	reader := bufio.NewReader(bytes.NewReader([]byte(input)))

	value, err := ReadRESP(reader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(value.Array) != 3 {
		t.Errorf("Expected array length 3, got %d", len(value.Array))
	}

	command, err := value.ToCommand()
	if err != nil {
		t.Fatalf("Unexpected error converting to command: %v", err)
	}

	if len(command) != 3 {
		t.Errorf("Expected command length 3, got %d", len(command))
	}

	if command[0] != "SET" {
		t.Errorf("Expected 'SET', got '%s'", command[0])
	}

	if command[1] != "mykey" {
		t.Errorf("Expected 'mykey', got '%s'", command[1])
	}

	if command[2] != "myvalue" {
		t.Errorf("Expected 'myvalue', got '%s'", command[2])
	}
}

func TestSerializeSimpleString(t *testing.T) {
	result := SerializeSimpleString("OK")
	expected := "+OK\r\n"

	if string(result) != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestSerializeError(t *testing.T) {
	result := SerializeError("ERR unknown command")
	expected := "-ERR unknown command\r\n"

	if string(result) != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestSerializeInteger(t *testing.T) {
	result := SerializeInteger(42)
	expected := ":42\r\n"

	if string(result) != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestSerializeBulkString(t *testing.T) {
	result := SerializeBulkString("hello")
	expected := "$5\r\nhello\r\n"

	if string(result) != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestSerializeBulkString_WithSpaces(t *testing.T) {
	result := SerializeBulkString("hello world")
	expected := "$11\r\nhello world\r\n"

	if string(result) != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestSerializeNullBulkString(t *testing.T) {
	result := SerializeNullBulkString()
	expected := "$-1\r\n"

	if string(result) != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestSerializeArray(t *testing.T) {
	elements := [][]byte{
		SerializeBulkString("hello"),
		SerializeBulkString("world"),
	}

	result := SerializeArray(elements)
	expected := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"

	if string(result) != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestRoundTrip_SimpleString(t *testing.T) {
	original := "PONG"
	serialized := SerializeSimpleString(original)

	reader := bufio.NewReader(bytes.NewReader(serialized))
	value, err := ReadRESP(reader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if value.Str != original {
		t.Errorf("Expected %q, got %q", original, value.Str)
	}
}

func TestRoundTrip_BulkString(t *testing.T) {
	original := "hello world"
	serialized := SerializeBulkString(original)

	reader := bufio.NewReader(bytes.NewReader(serialized))
	value, err := ReadRESP(reader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if value.Bulk != original {
		t.Errorf("Expected %q, got %q", original, value.Bulk)
	}
}

func TestRoundTrip_Integer(t *testing.T) {
	original := 42
	serialized := SerializeInteger(original)

	reader := bufio.NewReader(bytes.NewReader(serialized))
	value, err := ReadRESP(reader)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if value.Num != original {
		t.Errorf("Expected %d, got %d", original, value.Num)
	}
}

