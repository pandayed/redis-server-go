package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"strings"
	"sync"
)

// The flag package parses command-line arguments and converts them into Go variables. It's Go's built-in way to handle CLI options.
// The bufio package provides buffered I/O. It's useful for efficient reading and writing of data.

var (
	connectionsCount int
	connectionsMutex sync.Mutex
)

func main() {

	// host and port together determine where the server runs
	host := flag.String("host", "localhost", "Host to listen on")
	port := flag.String("port", "6379", "Port to listen on")

	help := flag.Bool("help", false, "Show help")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	address := *host + ":" + *port

	// Creates a listener on the specified address.
	listener, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatal("Failed to start server:", err)
	}

	log.Printf("Redis server listening on %s", address)

	storeInstance = newStore()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	connectionsMutex.Lock()
	connectionsCount++
	connectionsMutex.Unlock()

	log.Printf("New connection from %s. Total connections: %d", conn.RemoteAddr(), connectionsCount)

	defer func() {
		conn.Close()
		connectionsMutex.Lock()
		connectionsCount--
		connectionsMutex.Unlock()
		log.Printf("Connection from %s closed. Total connections: %d", conn.RemoteAddr(), connectionsCount)
	}()

	reader := bufio.NewReader(conn)

	for {
		value, err := ReadRESP(reader)
		if err != nil {
			if err.Error() == "EOF" {
				return
			}
			log.Printf("Error reading RESP: %v", err)
			return
		}

		command, err := value.ToCommand()
		if err != nil {
			response := SerializeError("ERR " + err.Error())
			conn.Write(response)
			continue
		}

		response := processCommand(command)
		conn.Write(response)
	}
}

func processCommand(command []string) []byte {
	if len(command) == 0 {
		return SerializeError("ERR empty command")
	}

	cmd := strings.ToUpper(command[0])

	switch cmd {

	case "PING":
		if len(command) == 1 {
			return SerializeSimpleString("PONG")
		}
		return SerializeBulkString(command[1])

	case "ECHO":
		if len(command) < 2 {
			return SerializeError("ERR wrong number of arguments for 'echo' command")
		}
		return SerializeBulkString(command[1])

	case "SET":
		if len(command) < 3 {
			return SerializeError("ERR wrong number of arguments for 'set' command")
		}
		storeInstance.Set(command[1], command[2])
		return SerializeSimpleString("OK")

	case "GET":
		if len(command) < 2 {
			return SerializeError("ERR wrong number of arguments for 'get' command")
		}
		value, exists := storeInstance.Get(command[1])
		if !exists {
			return SerializeNullBulkString()
		}
		return SerializeBulkString(value)

	case "INCR":
		if len(command) < 2 {
			return SerializeError("ERR wrong number of arguments for 'incr' command")
		}
		num, err := storeInstance.Incr(command[1])
		if err != nil {
			return SerializeError("ERR value is not an integer or out of range")
		}
		return SerializeInteger(num)

	case "DECR":
		if len(command) < 2 {
			return SerializeError("ERR wrong number of arguments for 'decr' command")
		}
		num, err := storeInstance.Decr(command[1])
		if err != nil {
			return SerializeError("ERR value is not an integer or out of range")
		}
		return SerializeInteger(num)

	case "DEL":
		if len(command) < 2 {
			return SerializeError("ERR wrong number of arguments for 'del' command")
		}
		deleted := storeInstance.Delete(command[1])
		if deleted {
			return SerializeInteger(1)
		}
		return SerializeInteger(0)

	case "EXISTS":
		if len(command) < 2 {
			return SerializeError("ERR wrong number of arguments for 'exists' command")
		}
		exists := storeInstance.Exists(command[1])
		if exists {
			return SerializeInteger(1)
		}
		return SerializeInteger(0)

	default:
		return SerializeError("ERR unknown command '" + cmd + "'")
	}

}
