package main

import (
	"testing"
)

func TestProcessCommand_PING(t *testing.T) {
	storeInstance = newStore()

	response := processCommand("PING")
	if response != "PONG" {
		t.Errorf("Expected PONG, got %s", response)
	}
}

func TestProcessCommand_ECHO(t *testing.T) {
	storeInstance = newStore()

	response := processCommand("ECHO hello world")
	if response != "hello world" {
		t.Errorf("Expected 'hello world', got %s", response)
	}
}

func TestProcessCommand_SET_GET(t *testing.T) {
	storeInstance = newStore()

	response := processCommand("SET mykey myvalue")
	if response != "OK" {
		t.Errorf("Expected OK, got %s", response)
	}

	response = processCommand("GET mykey")
	if response != "myvalue" {
		t.Errorf("Expected myvalue, got %s", response)
	}
}
