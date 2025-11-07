package main

import (
	"flag"
	"log"
	"net"
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

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		log.Printf("New connection from %s", conn.RemoteAddr())
	}
}
