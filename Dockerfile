# syntax=docker/dockerfile:1

# Build
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS build
WORKDIR /src
RUN apk add --no-cache build-base

ENV CGO_ENABLED=0
ENV GOFLAGS=-mod=vendor

COPY go.mod go.sum ./
COPY vendor/ ./vendor/
COPY . .

ARG TARGETOS TARGETARCH
ARG VERSION=dev
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath \
        -ldflags="-s -w -X main.version=$VERSION" \
        -o /out/blackmirror .

# Runtime
FROM alpine:3.22
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /out/blackmirror /app/blackmirror
EXPOSE 8080
ENTRYPOINT ["/app/blackmirror"]
