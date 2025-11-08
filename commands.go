package main

type CommandHandler func([]string) []byte

var commandRegistry = make(map[string]CommandHandler)

func registerCommand(name string, handler CommandHandler) {
	commandRegistry[name] = handler
}

func init() {
	registerConnectionCommands()
	registerStringCommands()
	registerListCommands()
	registerSetCommands()
	registerHashCommands()
}

func executeCommand(command []string) []byte {
	if len(command) == 0 {
		return SerializeError("ERR empty command")
	}

	cmdName := command[0]
	handler, exists := commandRegistry[cmdName]
	if !exists {
		return SerializeError("ERR unknown command '" + cmdName + "'")
	}

	return handler(command)
}

