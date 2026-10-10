package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"strings"

	"github.com/pandayed/redis-server-go/commands"
	"github.com/pandayed/redis-server-go/resp"
	"github.com/pandayed/redis-server-go/store"
)

func main() {
	host := flag.String("host", "localhost", "Host to listen on")
	port := flag.String("port", "6379", "Port to listen on")
	help := flag.Bool("help", false, "Show help")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	address := *host + ":" + *port

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}

	log.Printf("Redis server listening on %s", address)

	dataStore := store.New()
	executor := commands.New(dataStore)
	connManager := NewConnectionManager()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn, connManager, executor)
	}
}

func handleConnection(conn net.Conn, connManager *ConnectionManager, executor *commands.Executor) {
	connManager.Increment(conn.RemoteAddr())
	defer func() {
		conn.Close()
		connManager.Decrement(conn.RemoteAddr())
	}()

	reader := bufio.NewReader(conn)

	for {
		value, err := resp.ReadRESP(reader)
		if err != nil {
			if err.Error() == "EOF" {
				return
			}
			log.Printf("Error reading RESP: %v", err)
			return
		}

		command, err := value.ToCommand()
		if err != nil {
			response := resp.SerializeError("ERR " + err.Error())
			conn.Write(response)
			continue
		}

		cmdUpper := make([]string, len(command))
		for i, arg := range command {
			if i == 0 {
				cmdUpper[i] = strings.ToUpper(arg)
			} else {
				cmdUpper[i] = arg
			}
		}

		response := executor.Execute(cmdUpper)
		conn.Write(response)
	}
}
