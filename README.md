# trips

Trips is a microservice that calculates the beginning and ending airports for a series of flights.

Endpoints consume and produce JSON encoded in UTF-8.
Error information appears in an "error" result key.

## Use

Run:

```sh
$ go run main.go
```

## Endpoints

### Calculate

```
POST /calculate
{"data":{"flights":[["A","B"],["B","C"]]}}
=>
{"data":{"trip":["A","C"]}}
```

#### Example

Command:

```sh
$ curl -X POST -H 'Content-Type: application/json; charset=utf-8' -d '{"data":{"flights":[["A","B"],["B","C"]]}}' -s http://localhost:8080/calculate | jq
```

Result:

```
{
  "data": {
    "trip": [
      "A",
      "C"
    ]
  }
}
```
