package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ReadRESP(r *bufio.Reader) (string, error) {
	prefix, err := r.ReadByte()
	if err != nil {
		return "", err
	}

	switch prefix {
	case '+', '-', ':':
		// Simple string, error, integer — read until LF
		line, err := r.ReadString('\n')
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%c%s", prefix, line), nil

	case '$':
		// Bulk string: first read length
		lenLine, err := r.ReadString('\n')
		if err != nil {
			return "", err
		}

		length, err := strconv.Atoi(strings.TrimSpace(lenLine))
		if err != nil {
			return "", err
		}

		if length == -1 {
			return "(nil)\n", nil
		}

		// Read bulk data + CRLF
		buf := make([]byte, length+2)
		_, err = io.ReadFull(r, buf)
		if err != nil {
			return "", err
		}

		return string(buf[:length]) + "\n", nil

	default:
		return "", fmt.Errorf("unknown RESP prefix: %c", prefix)
	}
}

func ReadRESPArray(r *bufio.Reader) ([]string, error) {
	// First line: *<num>
	header, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}

	header = strings.TrimSpace(header)
	if !strings.HasPrefix(header, "*") {
		return nil, fmt.Errorf("invalid array header")
	}

	n, err := strconv.Atoi(header[1:])
	if err != nil {
		return nil, err
	}

	parts := make([]string, 0, n)
	for i := 0; i < n; i++ {
		// Read bulk length: $<len>
		bulkLenLine, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}

		bulkLenLine = strings.TrimSpace(bulkLenLine)
		if !strings.HasPrefix(bulkLenLine, "$") {
			return nil, fmt.Errorf("invalid bulk length")
		}

		bulkLen, err := strconv.Atoi(bulkLenLine[1:])
		if err != nil {
			return nil, err
		}

		buf := make([]byte, bulkLen+2) // include CRLF
		_, err = io.ReadFull(r, buf)
		if err != nil {
			return nil, err
		}

		parts = append(parts, string(buf[:bulkLen]))
	}

	return parts, nil
}
