# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/aodin/blackstar

# Build the command inside the container.
RUN go install github.com/aodin/blackstar

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/blackstar

# Document that the service listens on port 8081.
EXPOSE 8081
