package server

import (
	"bufio"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/stevenxchung/redis.go/internal/protocol"
)

type ValueWithExpiration struct {
	Value   string
	Expires *time.Time
}

type QueryHandler struct {
	inMemoryDB map[string]ValueWithExpiration
}

func NewQueryHandler() *QueryHandler {
	return &QueryHandler{
		inMemoryDB: make(map[string]ValueWithExpiration),
	}
}

func (qh *QueryHandler) queryHandler(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		// Parses a command from the reader
		command, err := reader.ReadString('\n')
		if err != nil {
			conn.Write([]byte(protocol.EncodeError("failed to parse command")))
			return
		}
		command = strings.TrimSpace(command)

		// Process command
		response := qh.processCommand(command)

		// Send response back to client
		_, err = writer.WriteString(response + "\n")
		if err != nil {
			conn.Write([]byte(protocol.EncodeError("failed to write response")))
			return
		}
		writer.Flush()
	}
}

func (qh *QueryHandler) processCommand(command string) string {
	input := strings.Fields(command)
	if len(input) == 0 {
		return protocol.EncodeError("empty command")
	}

	// Convert input command to upper case before checking
	switch strings.ToUpper(input[0]) {
	case "GET":
		return qh.get(input)
	case "SET":
		return qh.set(input)
	case "DEL":
		return qh.del(input)
	default:
		return protocol.EncodeError("unknown command: " + input[0])
	}
}

func (qh *QueryHandler) get(input []string) string {
	if len(input) != 2 {
		return protocol.EncodeError("wrong number of arguments for GET command")
	}

	key := input[1]
	object, found := qh.inMemoryDB[key]
	if !found {
		// Not found
		return protocol.NotFound()
	}

	if object.Expires != nil && time.Now().After(*object.Expires) {
		delete(qh.inMemoryDB, key)
		// Expired
		return protocol.NotFound()
	}

	return protocol.EncodeValue(object.Value)
}

func (qh *QueryHandler) set(input []string) string {
	if len(input) < 3 {
		return protocol.EncodeError("wrong number of arguments for SET command")
	}

	key := input[1]
	value := input[2]

	var (
		nx, xx       bool
		getLastValue bool
		ttl          *time.Time
	)

	i := 3
	for i < len(input) {
		arg := strings.ToUpper(input[i])
		switch arg {
		case "NX":
			// Set only if key does not exist
			nx = true
			i++
		case "XX":
			// Set only if key already exists
			xx = true
			i++
		case "GET":
			// Retrieve last value before update
			getLastValue = true
			i++
		case "EX":
			if i+1 >= len(input) {
				return protocol.EncodeError("syntax error: EX requires seconds")
			}
			seconds, err := strconv.Atoi(input[i+1])
			if err != nil || seconds <= 0 {
				return protocol.EncodeError("invalid expire time in SET command")
			}
			t := time.Now().Add(time.Duration(seconds) * time.Second)
			ttl = &t
			i += 2
		default:
			return protocol.EncodeError("syntax error near: " + arg)
		}
	}

	// Check if key found and not expired
	object, found := qh.inMemoryDB[key]
	if found && object.Expires != nil && time.Now().After(*object.Expires) {
		// Expired, treat as non-existent
		delete(qh.inMemoryDB, key)
		found = false
	}

	// NX/XX logic
	if nx && found {
		if getLastValue {
			if found {
				return protocol.EncodeValue(object.Value)
			} else {
				return protocol.NotFound()
			}
		}
		return protocol.NotFound()
	}
	if xx && !found {
		if getLastValue {
			return protocol.NotFound()
		}
		return protocol.NotFound()
	}

	// Store old value before overwrite
	var oldValue string
	var hadOld bool
	if getLastValue && found {
		oldValue = object.Value
		hadOld = true
	}

	// Set new value
	qh.inMemoryDB[key] = ValueWithExpiration{
		Value:   value,
		Expires: ttl,
	}

	if getLastValue {
		if hadOld {
			return protocol.EncodeValue(oldValue)
		}
		return protocol.NotFound()
	}
	return protocol.OK()
}

func (qh *QueryHandler) del(input []string) string {
	if len(input) < 2 {
		return protocol.EncodeError("wrong number of arguments for DEL command")
	}

	keys := input[1:]
	deletedCount := 0
	for _, key := range keys {
		if _, found := qh.inMemoryDB[key]; found {
			delete(qh.inMemoryDB, key)
			deletedCount++
		}
	}

	return protocol.EncodeInteger(deletedCount)
}
