package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"strconv"
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

	case "LPUSH":
		if len(command) < 3 {
			return SerializeError("ERR wrong number of arguments for 'lpush' command")
		}
		length := storeInstance.LPush(command[1], command[2:]...)
		return SerializeInteger(length)

	case "RPUSH":
		if len(command) < 3 {
			return SerializeError("ERR wrong number of arguments for 'rpush' command")
		}
		length := storeInstance.RPush(command[1], command[2:]...)
		return SerializeInteger(length)

	case "LPOP":
		if len(command) < 2 {
			return SerializeError("ERR wrong number of arguments for 'lpop' command")
		}
		value, exists := storeInstance.LPop(command[1])
		if !exists {
			return SerializeNullBulkString()
		}
		return SerializeBulkString(value)

	case "RPOP":
		if len(command) < 2 {
			return SerializeError("ERR wrong number of arguments for 'rpop' command")
		}
		value, exists := storeInstance.RPop(command[1])
		if !exists {
			return SerializeNullBulkString()
		}
		return SerializeBulkString(value)

	case "LRANGE":
		if len(command) < 4 {
			return SerializeError("ERR wrong number of arguments for 'lrange' command")
		}
		start, err1 := strconv.Atoi(command[2])
		stop, err2 := strconv.Atoi(command[3])
		if err1 != nil || err2 != nil {
			return SerializeError("ERR value is not an integer or out of range")
		}
		values := storeInstance.LRange(command[1], start, stop)
		elements := make([][]byte, len(values))
		for i, v := range values {
			elements[i] = SerializeBulkString(v)
		}
		return SerializeArray(elements)

	case "LLEN":
		if len(command) < 2 {
			return SerializeError("ERR wrong number of arguments for 'llen' command")
		}
		length := storeInstance.LLen(command[1])
		return SerializeInteger(length)

	case "SADD":
		if len(command) < 3 {
			return SerializeError("ERR wrong number of arguments for 'sadd' command")
		}
		added := storeInstance.SAdd(command[1], command[2:]...)
		return SerializeInteger(added)

	case "SMEMBERS":
		if len(command) < 2 {
			return SerializeError("ERR wrong number of arguments for 'smembers' command")
		}
		members := storeInstance.SMembers(command[1])
		elements := make([][]byte, len(members))
		for i, m := range members {
			elements[i] = SerializeBulkString(m)
		}
		return SerializeArray(elements)

	case "SISMEMBER":
		if len(command) < 3 {
			return SerializeError("ERR wrong number of arguments for 'sismember' command")
		}
		isMember := storeInstance.SIsMember(command[1], command[2])
		if isMember {
			return SerializeInteger(1)
		}
		return SerializeInteger(0)

	case "SREM":
		if len(command) < 3 {
			return SerializeError("ERR wrong number of arguments for 'srem' command")
		}
		removed := storeInstance.SRem(command[1], command[2:]...)
		return SerializeInteger(removed)

	case "SCARD":
		if len(command) < 2 {
			return SerializeError("ERR wrong number of arguments for 'scard' command")
		}
		cardinality := storeInstance.SCard(command[1])
		return SerializeInteger(cardinality)

	case "HSET":
		if len(command) < 4 {
			return SerializeError("ERR wrong number of arguments for 'hset' command")
		}
		added := storeInstance.HSet(command[1], command[2], command[3])
		return SerializeInteger(added)

	case "HGET":
		if len(command) < 3 {
			return SerializeError("ERR wrong number of arguments for 'hget' command")
		}
		value, exists := storeInstance.HGet(command[1], command[2])
		if !exists {
			return SerializeNullBulkString()
		}
		return SerializeBulkString(value)

	case "HGETALL":
		if len(command) < 2 {
			return SerializeError("ERR wrong number of arguments for 'hgetall' command")
		}
		hash := storeInstance.HGetAll(command[1])
		elements := make([][]byte, 0, len(hash)*2)
		for field, value := range hash {
			elements = append(elements, SerializeBulkString(field))
			elements = append(elements, SerializeBulkString(value))
		}
		return SerializeArray(elements)

	case "HDEL":
		if len(command) < 3 {
			return SerializeError("ERR wrong number of arguments for 'hdel' command")
		}
		deleted := storeInstance.HDel(command[1], command[2:]...)
		return SerializeInteger(deleted)

	case "HEXISTS":
		if len(command) < 3 {
			return SerializeError("ERR wrong number of arguments for 'hexists' command")
		}
		exists := storeInstance.HExists(command[1], command[2])
		if exists {
			return SerializeInteger(1)
		}
		return SerializeInteger(0)

	case "HLEN":
		if len(command) < 2 {
			return SerializeError("ERR wrong number of arguments for 'hlen' command")
		}
		length := storeInstance.HLen(command[1])
		return SerializeInteger(length)

	default:
		return SerializeError("ERR unknown command '" + cmd + "'")
	}

}
