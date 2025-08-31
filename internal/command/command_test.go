package command

import (
	"strings"
	"testing"
	"time"

	"github.com/stevenxchung/redis.go/internal/model"
	"github.com/stevenxchung/redis.go/internal/protocol"
)

func TestGet_Basic(t *testing.T) {
	db := make(map[string]model.ValueWithExpiration)
	db["foo"] = model.ValueWithExpiration{Value: "bar"}

	got := Get(db, []string{"GET", "foo"})
	want := protocol.EncodeValue("bar")
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestGet_KeyNotFound(t *testing.T) {
	db := make(map[string]model.ValueWithExpiration)
	got := Get(db, []string{"GET", "missing"})
	if got != protocol.NotFound() {
		t.Errorf("expected %q, got %q", protocol.NotFound(), got)
	}
}

func TestGet_ExpiredKey(t *testing.T) {
	db := make(map[string]model.ValueWithExpiration)
	expired := time.Now().Add(-time.Minute)
	db["k"] = model.ValueWithExpiration{Value: "v", Expires: &expired}

	got := Get(db, []string{"GET", "k"})
	if got != protocol.NotFound() {
		t.Errorf("expected not found for expired key, got %q", got)
	}
	if _, exists := db["k"]; exists {
		t.Error("expected expired key to be deleted from db")
	}
}

func TestDel_Basic(t *testing.T) {
	db := map[string]model.ValueWithExpiration{
		"a": {Value: "1"},
		"b": {Value: "2"},
	}
	got := Del(db, []string{"DEL", "a", "c", "b"})
	want := protocol.EncodeInteger(2)

	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
	if _, exists := db["a"]; exists {
		t.Error("expected a to be deleted")
	}
	if _, exists := db["b"]; exists {
		t.Error("expected b to be deleted")
	}
}

func TestDel_WrongArgs(t *testing.T) {
	db := make(map[string]model.ValueWithExpiration)
	got := Del(db, []string{"DEL"})
	want := protocol.EncodeError("wrong number of arguments for DEL command")
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestSet_Basic(t *testing.T) {
	db := make(map[string]model.ValueWithExpiration)
	got := Set(db, []string{"SET", "key1", "value1"})
	if got != protocol.OK() {
		t.Fatalf("expected OK, got %q", got)
	}
	if db["key1"].Value != "value1" {
		t.Errorf("expected value1, got %q", db["key1"].Value)
	}
}

func TestSet_WrongArgs(t *testing.T) {
	db := make(map[string]model.ValueWithExpiration)
	got := Set(db, []string{"SET", "keyOnly"})
	want := protocol.EncodeError("wrong number of arguments for SET command")
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestSet_NX_OnlyWhenNotExist(t *testing.T) {
	db := make(map[string]model.ValueWithExpiration)
	resp := Set(db, []string{"SET", "k", "v", "NX"})
	if resp != protocol.OK() {
		t.Fatalf("expected OK, got %q", resp)
	}

	// Try again with NX, should refuse overwrite
	resp = Set(db, []string{"SET", "k", "newv", "NX"})
	if resp != protocol.NotFound() {
		t.Errorf("expected not found on NX overwrite, got %q", resp)
	}

	if db["k"].Value != "v" {
		t.Errorf("value should not have been overwritten")
	}
}

func TestSet_XX_OnlyWhenExist(t *testing.T) {
	db := make(map[string]model.ValueWithExpiration)

	// Key doesn't exist yet, should fail
	resp := Set(db, []string{"SET", "k", "v", "XX"})
	if resp != protocol.NotFound() {
		t.Errorf("expected not found for XX on non-existent key, got %q", resp)
	}

	// Now it exists
	db["k"] = model.ValueWithExpiration{Value: "old"}
	resp = Set(db, []string{"SET", "k", "new", "XX"})
	if resp != protocol.OK() {
		t.Errorf("expected OK for XX on existing key, got %q", resp)
	}
	if db["k"].Value != "new" {
		t.Errorf("expected new, got %q", db["k"].Value)
	}
}

func TestSet_GET_ReturnsOldValue(t *testing.T) {
	db := make(map[string]model.ValueWithExpiration)

	// First time, no old value
	resp := Set(db, []string{"SET", "k", "v", "GET"})
	if resp != protocol.NotFound() {
		t.Errorf("expected not found on first set when using GET, got %q", resp)
	}

	// Now change it, should get old value
	db["k"] = model.ValueWithExpiration{Value: "old"}
	resp = Set(db, []string{"SET", "k", "new", "GET"})
	if resp != protocol.EncodeValue("old") {
		t.Errorf("expected old, got %q", resp)
	}
}

func TestSet_EX_Expiry(t *testing.T) {
	db := make(map[string]model.ValueWithExpiration)

	resp := Set(db, []string{"SET", "temp", "v", "EX", "1"})
	if resp != protocol.OK() {
		t.Fatalf("expected OK, got %q", resp)
	}

	// Expire after ~1s
	time.Sleep(1100 * time.Millisecond)

	getResp := Get(db, []string{"GET", "temp"})
	if getResp != protocol.NotFound() {
		t.Errorf("expected not found after expiry, got %q", getResp)
	}
}

func TestSet_InvalidEXArgs(t *testing.T) {
	db := make(map[string]model.ValueWithExpiration)

	resp := Set(db, []string{"SET", "k", "v", "EX"})
	if !strings.Contains(resp, "syntax error: EX requires seconds") {
		t.Errorf("expected EX requires seconds error, got %q", resp)
	}

	resp = Set(db, []string{"SET", "k", "v", "EX", "abc"})
	if !strings.Contains(resp, "invalid expire time") {
		t.Errorf("expected invalid expire time error, got %q", resp)
	}
}

func TestSet_ConflictingNXandXX(t *testing.T) {
	db := make(map[string]model.ValueWithExpiration)
	resp := Set(db, []string{"SET", "k", "v", "NX", "XX"})
	if !strings.Contains(resp, "not compatible") {
		t.Errorf("expected NX and XX conflict error, got %q", resp)
	}
}
