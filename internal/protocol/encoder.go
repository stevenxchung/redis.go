package protocol

import (
	"fmt"
	"strings"
)

func OK() string {
	return "+OK\r\n"
}

func NotFound() string {
	return "$-1\r\n"
}

func EncodeValue(value string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
}

func EncodeInteger(value int) string {
	return fmt.Sprintf(":%d\r\n", value)
}

func EncodeError(message string) string {
	return fmt.Sprintf("-ERR %s\r\n", message)
}

func EncodeRESPArray(args []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*%d\r\n", len(args)))
	for _, arg := range args {
		sb.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
	}
	return sb.String()
}
