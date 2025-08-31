package server

import (
	"bufio"
	"net"
	"testing"
	"time"

	"github.com/stevenxchung/redis.go/internal/model"
	"github.com/stevenxchung/redis.go/internal/protocol"
)

// ----- Unit Tests -----

func TestProcessCommand_EmptyCommand(t *testing.T) {
	qh := NewQueryHandler()
	got := qh.processCommand("")
	want := protocol.EncodeError("empty command")
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestProcessCommand_UnknownCommand(t *testing.T) {
	qh := NewQueryHandler()
	got := qh.processCommand("FOO")
	want := protocol.EncodeError("unknown command: FOO")
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestProcessCommand_SetAndGet(t *testing.T) {
	qh := NewQueryHandler()

	// Simple SET
	resp := qh.processCommand("SET mykey myval")
	if resp != protocol.OK() {
		t.Fatalf("expected OK, got %q", resp)
	}

	// GET
	got := qh.processCommand("GET mykey")
	want := protocol.EncodeValue("myval")
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestProcessCommand_GetNonExistent(t *testing.T) {
	qh := NewQueryHandler()
	got := qh.processCommand("GET missing")
	want := protocol.NotFound()
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestProcessCommand_Del(t *testing.T) {
	qh := NewQueryHandler()
	qh.inMemoryDB["k1"] = model.ValueWithExpiration{Value: "v1"}
	qh.inMemoryDB["k2"] = model.ValueWithExpiration{Value: "v2"}

	resp := qh.processCommand("DEL k1 k2 k3")
	// DEL returns integer count
	want := protocol.EncodeInteger(2)
	if resp != want {
		t.Errorf("expected %q, got %q", want, resp)
	}
	if _, exists := qh.inMemoryDB["k1"]; exists {
		t.Errorf("k1 should have been deleted")
	}
	if _, exists := qh.inMemoryDB["k2"]; exists {
		t.Errorf("k2 should have been deleted")
	}
}

func TestProcessCommand_SetWithExpiry(t *testing.T) {
	qh := NewQueryHandler()

	// EX 1 = expire after 1 second
	resp := qh.processCommand("SET tempkey tempval EX 1")
	if resp != protocol.OK() {
		t.Fatalf("expected OK, got %q", resp)
	}

	time.Sleep(1500 * time.Millisecond) // allow expiry

	// GET should now return not found
	got := qh.processCommand("GET tempkey")
	want := protocol.NotFound()
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

// ----- E2E Tests -----

func TestQueryHandler_TCPFlow(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	qh := NewQueryHandler()
	go qh.queryHandler(serverConn)

	clientWriter := bufio.NewWriter(clientConn)
	clientReader := bufio.NewReader(clientConn)

	sendRESP := func(args ...string) string {
		_, _ = clientWriter.WriteString(protocol.EncodeRESPArray(args))
		clientWriter.Flush()
		resp, err := protocol.ReadRESP(clientReader) // decoded result
		if err != nil {
			t.Fatalf("failed to read RESP: %v", err)
		}
		return resp
	}

	// SET should decode to "+OK\r\n"
	if got := sendRESP("SET", "foo", "bar"); got != "+OK\r\n" {
		t.Fatalf("expected %q, got %q", "+OK\r\n", got)
	}

	// GET should decode to "bar\n"
	if got := sendRESP("GET", "foo"); got != "bar\n" {
		t.Fatalf("expected %q, got %q", "bar\n", got)
	}

	// DEL should decode to ":1\r\n"
	if got := sendRESP("DEL", "foo"); got != ":1\r\n" {
		t.Fatalf("expected %q, got %q", ":1\r\n", got)
	}

	clientConn.Close()
}
