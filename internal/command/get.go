package command

import (
	"time"

	"github.com/stevenxchung/redis.go/internal/model"
	"github.com/stevenxchung/redis.go/internal/protocol"
)

func Get(db map[string]model.ValueWithExpiration, input []string) string {
	if len(input) != 2 {
		return protocol.EncodeError("wrong number of arguments for GET command")
	}

	key := input[1]
	object, found := db[key]
	if !found {
		// Not found
		return protocol.NotFound()
	}

	if object.Expires != nil && time.Now().After(*object.Expires) {
		delete(db, key)
		// Expired
		return protocol.NotFound()
	}

	return protocol.EncodeValue(object.Value)
}
