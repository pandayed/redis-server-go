package main

import "github.com/pandayed/redis-server-go/commands"

// executeCommand preserves the entry point used by existing tests.
func executeCommand(command []string) []byte {
	return commands.New(storeInstance).Execute(command)
}
