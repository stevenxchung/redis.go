# redis.go

[Redis](https://redis.io/) is a well-known in-memory database that persists data on disk and operates on a key-value data model. To better understand its internals, we can explore building a similar solution from scratch. This repository provides such an implementation, incorporating some variations.

We will start with a simple GraphQL API interface where users can make simple queries (e.g., via Postman) to this redis.go database and go from there.

## Setup

To run locally:

1. Start the app: `go run main.go`
2. Once message occurs: `[INFO] Server started on http://localhost:8000` use the address to configure your requests in Postman
3. Make a request: use Postman or any other HTTP client to send requests to the server

**Optional**: create a `.env` file with `REDIS_GO_PORT=<YOUR PORT>` to adjust the port

## GraphQL API Reference

Note that the most common convention for GraphQL requests is HTTP [POST](https://graphql.org/learn/serving-over-http/#http-methods-headers-and-body) and `Content-Type: application/json`.

### GET

Get the value of key. If the key does not exist or has expired, the service will provide a message.

```graphql
query {
  get(key: "test")
}
```

### SET

Sets the value of a key. Expiration time is an optional parameter.

```graphql
mutation {
  set(key: "test", value: "Over 9000!", expires: 3)
}
```

### DEL

Removes the specified key(s). A key is ignored if it does not exist.

```graphql
mutation {
  del(keys: ["test"])
}
```
