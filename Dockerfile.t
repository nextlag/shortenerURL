FROM golang:1.21.1-alpine

# Install git
RUN set -ex; \
    apk update; \
    apk add --no-cache git

# Set working directory
WORKDIR /go/src/github.com/nextbug/shortenerURL
COPY . .

# Run tests
CMD CGO_ENABLED=0 go test -count=1 -v ./...