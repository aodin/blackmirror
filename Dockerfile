# Start from an Alpine Linux image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.10-alpine3.7 AS build

# Copy the local package files to the container's workspace.
COPY . /go/src/github.com/aodin/blackmirror

# Build the command inside the container. Since we're using Alpine Linux
# to both build and run, we can keep CGO enabled.
# Note that Alpine Linux uses the musl C library (https://www.musl-libc.org/)
# and not the GNU C Library used by most other Linux distros.
RUN go install github.com/aodin/blackmirror

# The default Alpine Linux images is about 4 MB.
FROM alpine:3.7 AS run
LABEL description="https://github.com/aodin/blackmirror"

# Alpine does not ship with root CA certificates. We'll install them using
# the --no-cache flag, since the image's package index may be stale.
RUN apk add --no-cache ca-certificates

# Copy the binary from the build stage.
COPY --from=build /go/bin/blackmirror .

# Run the command by default when the container starts.
# Using exec format runs the Go server at PID 1, allowing "docker stop"
# SIGTERM to gracefully stop the Go server.
ENTRYPOINT ["./blackmirror"]

# Blackmirror's default port is 8080, be can be overridden with either
# The PORT env var or -port flag
EXPOSE 8080
