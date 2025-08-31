package command

import (
	"strconv"
	"strings"
	"time"

	"github.com/stevenxchung/redis.go/internal/model"
	"github.com/stevenxchung/redis.go/internal/protocol"
)

type SetOptions struct {
	NX, XX, GET bool
	TTL         *time.Time
}

func Set(db map[string]model.ValueWithExpiration, input []string) string {
	if len(input) < 3 {
		return protocol.EncodeError("wrong number of arguments for SET command")
	}

	key := input[1]
	value := input[2]

	opts, errMsg := parseSetOptions(input[3:])
	if errMsg != "" {
		return errMsg
	}

	if opts.NX && opts.XX {
		return protocol.EncodeError("syntax error: NX and XX options at the same time are not compatible")
	}

	// Expiration check: remove expired keys before logic
	object, found := fetchAndRemoveIfExpired(db, key)

	// Apply NX/XX constraints
	if violatesExistenceRules(opts, found) {
		return existingValueOrNull(opts, found, object)
	}

	// Save old value if GET requested
	oldValue, hadOld := "", false
	if opts.GET && found {
		oldValue = object.Value
		hadOld = true
	}

	// Set the new value
	db[key] = model.ValueWithExpiration{
		Value:   value,
		Expires: opts.TTL,
	}

	// Return according to GET logic
	if opts.GET {
		if hadOld {
			return protocol.EncodeValue(oldValue)
		}
		return protocol.NotFound()
	}

	return protocol.OK()
}

func parseSetOptions(args []string) (SetOptions, string) {
	opts := SetOptions{}
	i := 0
	for i < len(args) {
		switch strings.ToUpper(args[i]) {
		case "NX":
			// Set only if key does not exist
			opts.NX = true
			i++
		case "XX":
			// Set only if key already exists
			opts.XX = true
			i++
		case "GET":
			// Retrieve last value before update
			opts.GET = true
			i++
		case "EX":
			if i+1 >= len(args) {
				return opts, protocol.EncodeError("syntax error: EX requires seconds")
			}
			seconds, err := strconv.Atoi(args[i+1])
			if err != nil || seconds <= 0 {
				return opts, protocol.EncodeError("invalid expire time in SET command")
			}
			t := time.Now().Add(time.Duration(seconds) * time.Second)
			opts.TTL = &t
			i += 2
		default:
			return opts, protocol.EncodeError("syntax error near: " + args[i])
		}
	}
	return opts, ""
}

func fetchAndRemoveIfExpired(db map[string]model.ValueWithExpiration, key string) (model.ValueWithExpiration, bool) {
	obj, found := db[key]
	if found && obj.Expires != nil && time.Now().After(*obj.Expires) {
		delete(db, key)
		return model.ValueWithExpiration{}, false
	}
	return obj, found
}

func violatesExistenceRules(opts SetOptions, found bool) bool {
	if opts.NX && found {
		return true
	}
	if opts.XX && !found {
		return true
	}
	return false
}

func existingValueOrNull(opts SetOptions, found bool, obj model.ValueWithExpiration) string {
	if opts.GET && found {
		return protocol.EncodeValue(obj.Value)
	}
	return protocol.NotFound()
}
