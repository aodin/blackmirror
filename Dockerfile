# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/aodin/blackmirror

# Build the command inside the container.
RUN go install github.com/aodin/blackmirror

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/blackmirror

# Document that the service listens on port 8080.
EXPOSE 8080

# Or use onbuild
# FROM golang:onbuild
# EXPOSE 8080
