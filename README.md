# Redis Query | [Docker Hub](https://hub.docker.com/r/reeganexe/redis-query)

A simple cli/web utility that is used to extract values from Redis.

```
go get github.com/ReeganExE/redis-query
```

## Usage

```sh
redis-query -help

  -addr string
        Redis address (default "127.0.0.1:6379")
  -key string
        Key name

  -pattern string
        A regex pattern to extract the Value. e.g: -pattern "(\d+)"

  -port uint (Optional)
        HTTP port, if specified, app will start an HTTP server on the specified port

        GET http://0.0.0.0:port/query?key=your-key&pattern=(\d+)

        HTTP/1.1 200 OK
        Content-type: application/json
        Content-Length: 27
        Connection: Closed

        {"value": "extracted value"}

  -verbose
        Print raw Value
```

## Docker

```sh
docker run --rm \
  -p 8910:8910 \
  reeganexe/redis-query -addr 10.11.12.13:6379 -port 8910
```

Then,

```sh
curl -i "http://0.0.0.0:8910/query?key=your:key&pattern=It has been (\d+) years"

HTTP/1.1 200 OK
Content-type: application/json
Content-Length: 14
Connection: Closed

{"value": "3"}
```
