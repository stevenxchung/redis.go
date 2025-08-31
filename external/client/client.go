package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/stevenxchung/redis.go/internal/config"
	"github.com/stevenxchung/redis.go/internal/protocol"
	"github.com/stevenxchung/redis.go/pkg/util"
)

func StartClient(cfg *config.Config) {
	conn, err := net.Dial("tcp", ":"+cfg.ServerPort)
	if err != nil {
		util.LogError("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	util.LogInfo("Connecting to redis.go server...")
	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	for {
		fmt.Print("redis.go> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			util.LogInfo(fmt.Sprintf("Error reading command: %s", err))
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.EqualFold(line, "quit") || strings.EqualFold(line, "exit") {
			fmt.Println("Goodbye!")
			break
		}

		// Encode to RESP Array
		args := strings.Fields(line)
		resp := protocol.EncodeRESPArray(args)

		_, err = fmt.Fprint(conn, resp)
		if err != nil {
			util.LogInfo(fmt.Sprintf("Error sending command: %s", err))
			continue
		}

		// Parse full RESP reply
		response, err := protocol.ReadRESP(serverReader)
		if err != nil {
			util.LogInfo(fmt.Sprintf("Error reading response: %s", err))
			continue
		}

		fmt.Print(response)
	}
}
