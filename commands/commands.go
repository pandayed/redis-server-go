package commands

import (
	"github.com/pandayed/redis-server-go/resp"
	"github.com/pandayed/redis-server-go/store"
)

type CommandHandler func([]string) []byte

// Executor dispatches commands using the store supplied at construction.
type Executor struct {
	dataStore       store.DataStore
	commandRegistry map[string]CommandHandler
}

// New creates an executor and registers its supported commands.
func New(dataStore store.DataStore) *Executor {
	e := &Executor{
		dataStore:       dataStore,
		commandRegistry: make(map[string]CommandHandler),
	}
	e.registerConnectionCommands()
	e.registerStringCommands()
	e.registerListCommands()
	e.registerSetCommands()
	e.registerHashCommands()
	return e
}

func (e *Executor) registerCommand(name string, handler CommandHandler) {
	e.commandRegistry[name] = handler
}

// Execute dispatches a command whose name has already been uppercased.
func (e *Executor) Execute(command []string) []byte {
	if len(command) == 0 {
		return resp.SerializeError("ERR empty command")
	}

	cmdName := command[0]
	handler, exists := e.commandRegistry[cmdName]
	if !exists {
		return resp.SerializeError("ERR unknown command '" + cmdName + "'")
	}

	return handler(command)
}
