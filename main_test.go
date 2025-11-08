package main

import (
	"strings"
	"testing"
)

func executeTestCommand(command []string) []byte {
	cmdUpper := make([]string, len(command))
	for i, arg := range command {
		if i == 0 {
			cmdUpper[i] = strings.ToUpper(arg)
		} else {
			cmdUpper[i] = arg
		}
	}
	return executeCommand(cmdUpper)
}

func TestProcessCommand_PING(t *testing.T) {
	storeInstance = newStore()

	response := executeTestCommand([]string{"PING"})
	expected := SerializeSimpleString("PONG")

	if string(response) != string(expected) {
		t.Errorf("Expected %q, got %q", expected, response)
	}
}

func TestProcessCommand_PING_WithMessage(t *testing.T) {
	storeInstance = newStore()

	response := executeTestCommand([]string{"PING", "hello"})
	expected := SerializeBulkString("hello")

	if string(response) != string(expected) {
		t.Errorf("Expected %q, got %q", expected, response)
	}
}

func TestProcessCommand_ECHO(t *testing.T) {
	storeInstance = newStore()

	response := executeTestCommand([]string{"ECHO", "hello world"})
	expected := SerializeBulkString("hello world")

	if string(response) != string(expected) {
		t.Errorf("Expected %q, got %q", expected, response)
	}
}

func TestProcessCommand_SET_GET(t *testing.T) {
	storeInstance = newStore()

	response := executeTestCommand([]string{"SET", "mykey", "myvalue"})
	expected := SerializeSimpleString("OK")

	if string(response) != string(expected) {
		t.Errorf("Expected %q, got %q", expected, response)
	}

	response = executeTestCommand([]string{"GET", "mykey"})
	expected = SerializeBulkString("myvalue")

	if string(response) != string(expected) {
		t.Errorf("Expected %q, got %q", expected, response)
	}
}

func TestProcessCommand_GET_NonExistent(t *testing.T) {
	storeInstance = newStore()

	response := executeTestCommand([]string{"GET", "nonexistent"})
	expected := SerializeNullBulkString()

	if string(response) != string(expected) {
		t.Errorf("Expected %q, got %q", expected, response)
	}
}

func TestProcessCommand_INCR(t *testing.T) {
	storeInstance = newStore()

	response := executeTestCommand([]string{"INCR", "counter"})
	expected := SerializeInteger(1)

	if string(response) != string(expected) {
		t.Errorf("Expected %q, got %q", expected, response)
	}

	response = executeTestCommand([]string{"INCR", "counter"})
	expected = SerializeInteger(2)

	if string(response) != string(expected) {
		t.Errorf("Expected %q, got %q", expected, response)
	}
}

func TestProcessCommand_DECR(t *testing.T) {
	storeInstance = newStore()

	response := executeTestCommand([]string{"DECR", "counter"})
	expected := SerializeInteger(-1)

	if string(response) != string(expected) {
		t.Errorf("Expected %q, got %q", expected, response)
	}
}

func TestProcessCommand_DEL(t *testing.T) {
	storeInstance = newStore()

	storeInstance.Set("key1", "value1")

	response := executeTestCommand([]string{"DEL", "key1"})
	expected := SerializeInteger(1)

	if string(response) != string(expected) {
		t.Errorf("Expected %q, got %q", expected, response)
	}

	response = executeTestCommand([]string{"DEL", "nonexistent"})
	expected = SerializeInteger(0)

	if string(response) != string(expected) {
		t.Errorf("Expected %q, got %q", expected, response)
	}
}

func TestProcessCommand_EXISTS(t *testing.T) {
	storeInstance = newStore()

	storeInstance.Set("key1", "value1")

	response := executeTestCommand([]string{"EXISTS", "key1"})
	expected := SerializeInteger(1)

	if string(response) != string(expected) {
		t.Errorf("Expected %q, got %q", expected, response)
	}

	response = executeTestCommand([]string{"EXISTS", "nonexistent"})
	expected = SerializeInteger(0)

	if string(response) != string(expected) {
		t.Errorf("Expected %q, got %q", expected, response)
	}
}

func TestProcessCommand_CaseInsensitive(t *testing.T) {
	storeInstance = newStore()

	response := executeTestCommand([]string{"ping"})
	expected := SerializeSimpleString("PONG")

	if string(response) != string(expected) {
		t.Errorf("Expected %q, got %q", expected, response)
	}

	response = executeTestCommand([]string{"PiNg"})
	expected = SerializeSimpleString("PONG")

	if string(response) != string(expected) {
		t.Errorf("Expected %q, got %q", expected, response)
	}
}

func TestProcessCommand_UnknownCommand(t *testing.T) {
	storeInstance = newStore()

	response := executeTestCommand([]string{"UNKNOWN"})
	expected := SerializeError("ERR unknown command 'UNKNOWN'")

	if string(response) != string(expected) {
		t.Errorf("Expected %q, got %q", expected, response)
	}
}
