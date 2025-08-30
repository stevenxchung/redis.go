package command

import (
	"github.com/stevenxchung/redis.go/internal/model"
	"github.com/stevenxchung/redis.go/internal/protocol"
)

func Del(db map[string]model.ValueWithExpiration, input []string) string {
	if len(input) < 2 {
		return protocol.EncodeError("wrong number of arguments for DEL command")
	}

	keys := input[1:]
	deletedCount := 0
	for _, key := range keys {
		if _, found := db[key]; found {
			delete(db, key)
			deletedCount++
		}
	}

	return protocol.EncodeInteger(deletedCount)
}
