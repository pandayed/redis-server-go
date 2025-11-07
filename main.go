package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"strconv"
	"strings"
)

// The flag package parses command-line arguments and converts them into Go variables. It's Go's built-in way to handle CLI options.

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
		log.Printf("New connection from %s", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		line := scanner.Text()
		response := processCommand(line)
		conn.Write([]byte(response + "\r\n"))
	}

}

func processCommand(command string) string {
	parts := strings.Split(command, " ")

	if len(parts) == 0 {
		return "-ERR empty command"
	}

	switch parts[0] {

	case "PING":
		return "PONG"

	case "ECHO":
		return strings.Join(parts[1:], " ")

	case "SET":
		storeInstance.Set(parts[1], parts[2])
		return "OK"

	case "GET":
		value, exists := storeInstance.Get(parts[1])
		if !exists {
			return ""
		}
		return value

	case "INCR":
		num, err := storeInstance.Incr(parts[1])
		if err != nil {
			return "-ERR invalid increment"
		}
		return strconv.Itoa(num)

	case "DECR":
		num, err := storeInstance.Decr(parts[1])
		if err != nil {
			return "-ERR invalid decrement"
		}
		return strconv.Itoa(num)

	case "DEL":
		deleted := storeInstance.Delete(parts[1])
		if !deleted {
			return "FALSE"
		}
		return "TRUE"

	case "EXISTS":
		exists := storeInstance.Exists(parts[1])
		if !exists {
			return "FALSE"
		}
		return "TRUE"

	default:
		return "-ERR unknown command"
	}

}
