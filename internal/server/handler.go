package server

import (
	"bufio"
	"net"
	"strings"

	"github.com/stevenxchung/redis.go/internal/command"
	"github.com/stevenxchung/redis.go/internal/model"
	"github.com/stevenxchung/redis.go/internal/protocol"
)

type QueryHandler struct {
	inMemoryDB map[string]model.ValueWithExpiration
}

func NewQueryHandler() *QueryHandler {
	return &QueryHandler{
		inMemoryDB: make(map[string]model.ValueWithExpiration),
	}
}

func (qh *QueryHandler) queryHandler(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		parts, err := protocol.ReadRESPArray(reader)
		if err != nil {
			return
		}

		response := qh.processCommand(strings.Join(parts, " "))
		writer.WriteString(response)
		writer.Flush()
	}
}

func (qh *QueryHandler) processCommand(cmd string) string {
	input := strings.Fields(cmd)
	if len(input) == 0 {
		return protocol.EncodeError("empty command")
	}

	// Convert input command to upper case before checking
	switch strings.ToUpper(input[0]) {
	case "GET":
		return command.Get(qh.inMemoryDB, input)
	case "SET":
		return command.Set(qh.inMemoryDB, input)
	case "DEL":
		return command.Del(qh.inMemoryDB, input)
	default:
		return protocol.EncodeError("unknown command: " + input[0])
	}
}
