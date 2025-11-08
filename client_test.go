package main

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func TestRealConnection_PING(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Skip("Server not running, skipping integration test")
		return
	}
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	command := "*1\r\n$4\r\nPING\r\n"
	writer.WriteString(command)
	writer.Flush()

	response, err := ReadRESP(reader)
	if err != nil {
		t.Fatalf("Error reading response: %v", err)
	}

	if response.Type != SimpleString {
		t.Errorf("Expected SimpleString, got %c", response.Type)
	}

	if response.Str != "PONG" {
		t.Errorf("Expected PONG, got %s", response.Str)
	}
}

func TestRealConnection_SET_GET(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Skip("Server not running, skipping integration test")
		return
	}
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	setCommand := "*3\r\n$3\r\nSET\r\n$7\r\ntestkey\r\n$9\r\ntestvalue\r\n"
	writer.WriteString(setCommand)
	writer.Flush()

	response, err := ReadRESP(reader)
	if err != nil {
		t.Fatalf("Error reading SET response: %v", err)
	}

	if response.Type != SimpleString || response.Str != "OK" {
		t.Errorf("Expected OK, got %v", response)
	}

	getCommand := "*2\r\n$3\r\nGET\r\n$7\r\ntestkey\r\n"
	writer.WriteString(getCommand)
	writer.Flush()

	response, err = ReadRESP(reader)
	if err != nil {
		t.Fatalf("Error reading GET response: %v", err)
	}

	if response.Type != BulkString {
		t.Errorf("Expected BulkString, got %c", response.Type)
	}

	if response.Bulk != "testvalue" {
		t.Errorf("Expected 'testvalue', got '%s'", response.Bulk)
	}
}

func TestRealConnection_INCR(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Skip("Server not running, skipping integration test")
		return
	}
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	delCommand := "*2\r\n$3\r\nDEL\r\n$12\r\ntest_counter\r\n"
	writer.WriteString(delCommand)
	writer.Flush()
	ReadRESP(reader)

	incrCommand := "*2\r\n$4\r\nINCR\r\n$12\r\ntest_counter\r\n"
	writer.WriteString(incrCommand)
	writer.Flush()

	response, err := ReadRESP(reader)
	if err != nil {
		t.Fatalf("Error reading INCR response: %v", err)
	}

	if response.Type != Integer {
		t.Errorf("Expected Integer, got %c", response.Type)
	}

	if response.Num != 1 {
		t.Errorf("Expected 1, got %d", response.Num)
	}

	writer.WriteString(incrCommand)
	writer.Flush()

	response, err = ReadRESP(reader)
	if err != nil {
		t.Fatalf("Error reading second INCR response: %v", err)
	}

	if response.Num != 2 {
		t.Errorf("Expected 2, got %d", response.Num)
	}
}

func TestRealConnection_MultipleCommands(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Skip("Server not running, skipping integration test")
		return
	}
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	commands := []string{
		"*1\r\n$4\r\nPING\r\n",
		"*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n",
		"*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n",
		"*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n",
	}

	expectedResponses := []struct {
		respType RESPType
		value    string
	}{
		{SimpleString, "PONG"},
		{BulkString, "hello"},
		{SimpleString, "OK"},
		{BulkString, "value"},
	}

	for i, cmd := range commands {
		writer.WriteString(cmd)
		writer.Flush()

		response, err := ReadRESP(reader)
		if err != nil {
			t.Fatalf("Error reading response %d: %v", i, err)
		}

		expected := expectedResponses[i]
		if response.Type != expected.respType {
			t.Errorf("Command %d: Expected type %c, got %c", i, expected.respType, response.Type)
		}

		var actualValue string
		if response.Type == SimpleString {
			actualValue = response.Str
		} else if response.Type == BulkString {
			actualValue = response.Bulk
		}

		if actualValue != expected.value {
			t.Errorf("Command %d: Expected '%s', got '%s'", i, expected.value, actualValue)
		}
	}
}

func TestRealConnection_CaseInsensitive(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		t.Skip("Server not running, skipping integration test")
		return
	}
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	commands := []string{
		"*1\r\n$4\r\nping\r\n",
		"*1\r\n$4\r\nPING\r\n",
		"*1\r\n$4\r\nPiNg\r\n",
	}

	for i, cmd := range commands {
		writer.WriteString(cmd)
		writer.Flush()

		response, err := ReadRESP(reader)
		if err != nil {
			t.Fatalf("Error reading response %d: %v", i, err)
		}

		if response.Type != SimpleString || response.Str != "PONG" {
			t.Errorf("Command %d: Expected PONG, got %v", i, response)
		}
	}
}

func TestRealConnection_Concurrent(t *testing.T) {
	numClients := 10
	done := make(chan bool, numClients)

	for i := 0; i < numClients; i++ {
		go func(clientID int) {
			conn, err := net.Dial("tcp", "localhost:6379")
			if err != nil {
				t.Logf("Client %d: Could not connect: %v", clientID, err)
				done <- false
				return
			}
			defer conn.Close()

			writer := bufio.NewWriter(conn)
			reader := bufio.NewReader(conn)

			command := "*1\r\n$4\r\nPING\r\n"
			writer.WriteString(command)
			writer.Flush()

			response, err := ReadRESP(reader)
			if err != nil {
				t.Logf("Client %d: Error reading response: %v", clientID, err)
				done <- false
				return
			}

			if response.Type != SimpleString || response.Str != "PONG" {
				t.Logf("Client %d: Expected PONG, got %v", clientID, response)
				done <- false
				return
			}

			done <- true
		}(i)
	}

	timeout := time.After(5 * time.Second)
	successCount := 0

	for i := 0; i < numClients; i++ {
		select {
		case success := <-done:
			if success {
				successCount++
			}
		case <-timeout:
			t.Fatal("Test timed out")
		}
	}

	if successCount != numClients {
		t.Errorf("Expected %d successful clients, got %d", numClients, successCount)
	}
}
