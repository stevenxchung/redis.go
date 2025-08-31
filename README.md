# redis.go

[Redis](https://redis.io/) is a well-known in-memory database that persists data on disk and operates on a key-value data model. To better understand its internals, we can explore building a similar solution from scratch. This repository provides such an implementation, incorporating some variations.

We will start with a simple [Redis-based interface](https://redis.io/docs/latest/develop/reference/protocol-spec/) where users can make simple queries to this redis.go database via a built-in redis.go client and go from there.

## Setup

To run locally:

1. Start the app in a terminal: `go run main.go`
2. Once this message appears: `[INFO] Server started on localhost:6379` you may need to allow run permissions
3. A built-in client automatically connects to the server above: `[INFO] Connecting to redis.go server...`
4. Once the client prompt appears: `redis.go>` you may execute using Redis commands as shown below

**Optional**: create a `.env` file with `REDIS_GO_SERVER_PORT=<YOUR PORT>` to adjust the port

## redis.go API Reference

Commands are based on the official [Redis commands](https://redis.io/docs/latest/commands/). We will note some base commands below and progressively add more functionality.

### GET

Get the value of key. If the key does not exist the special value `nil` is returned.

```bash
GET key
```

### SET

Set key to hold the string value. If key already holds a value, it is overwritten. Any previous time to live associated with the key is discarded on successful SET operation.

```bash
SET key value [NX | XX] [GET] [EX seconds]
```

Where:

- `NX`: set only if key does **NOT** exist
- `XX`: set only if key **ALREADY** exists
- `GET`: retrieve last value before update

### DEL

Removes the specified keys. A key is ignored if it does not exist.

```bash
DEL key [key ...]
```

## Testing

```
# Run all tests in all subpackages and create a coverage profile
go test ./... -coverprofile=coverage.out

# View a simple coverage summary in the terminal
go tool cover -func=coverage.out

# Open an HTML coverage report in your browser
go tool cover -html=coverage.out
```
