package protocol

import (
	"fmt"
)

func OK() string {
	return "OK"
}

func NotFound() string {
	return "(nil)"
}

func EncodeValue(value string) string {
	return fmt.Sprintf("\"%s\"\r\n", value)
}

func EncodeInteger(value int) string {
	return fmt.Sprintf("(integer) %v\r\n", value)
}

func EncodeError(value string) string {
	return fmt.Sprintf("(error) %s\r\n", value)
}
