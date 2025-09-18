# Black Mirror

A small Go server that reflects request headers, useful for testing deployments.

### Install

A Docker image is available on [Docker Hub](https://hub.docker.com/r/aodin/blackmirror/).

    docker pull aodin/blackmirror
    docker run -p 8080:8080 aodin/blackmirror

It was built for multi-arch and tagged with:

    docker buildx create --use --name multi || true
    docker buildx inspect --bootstrap
    docker buildx build \
        --platform linux/amd64,linux/arm64 \
        -t aodin/blackmirror:latest \
        --push .

The server can be built locally with [go](https://golang.org/):

    go build -mod=vendor .


### Usage

The server will respond to a request for any method and path with the HTTP/1.1 wire format as a text/plain document.

Request:

    GET localhost:8080/yo

Response:

    GET /yo HTTP/1.1
    Host: localhost:8080
    Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
    Accept-Encoding: gzip, deflate
    Accept-Language: en-US,en;q=0.9
    Connection: keep-alive
    Priority: u=0, i
    Sec-Fetch-Dest: document
    Sec-Fetch-Mode: navigate
    Sec-Fetch-Site: none
    Upgrade-Insecure-Requests: 1
    User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.6 Safari/605.1.15

It accepts `PORT` / `HOST` environment variables and `-port` / `-host` flags, with the flags taking precendence. Its default address is `:8080`.

    PORT=8081 HOST=0.0.0.0 go run blackmirror.go -port=8082 -host="localhost"

Happy hacking!

aodin, 2025
