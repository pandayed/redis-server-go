package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"strings"
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

	storeInstance = newStore()
	connManager := NewConnectionManager()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn, connManager)
	}
}

func handleConnection(conn net.Conn, connManager *ConnectionManager) {
	connManager.Increment(conn.RemoteAddr())
	defer func() {
		conn.Close()
		connManager.Decrement(conn.RemoteAddr())
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

		cmdUpper := make([]string, len(command))
		for i, arg := range command {
			if i == 0 {
				cmdUpper[i] = strings.ToUpper(arg)
			} else {
				cmdUpper[i] = arg
			}
		}

		response := executeCommand(cmdUpper)
		conn.Write(response)
	}
}
