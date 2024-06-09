package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/stevenxchung/redis.go/internal/config"
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
	for {
		fmt.Print("redis.go> ")
		command, err := reader.ReadString('\n')
		if err != nil {
			util.LogInfo(fmt.Sprintf("Error reading command: %s", err))
			continue
		}
		command = strings.TrimSpace(command)

		if command == "quit" || command == "exit" {
			fmt.Print("Goodbye!")
			break
		}

		_, err = fmt.Fprintf(conn, "%s\n", command)
		if err != nil {
			util.LogInfo(fmt.Sprintf("Error sending command: %s", err))
			continue
		}

		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			util.LogInfo(fmt.Sprintf("Error reading response: %s", err))
			continue
		}

		fmt.Print(response)
	}
}
