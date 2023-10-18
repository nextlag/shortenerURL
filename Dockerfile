FROM golang:1.21.1-alpine AS src

# Install git
RUN set -ex; \
    apk update; \
    apk add --no-cache git

# Copy Repository
WORKDIR /go/src/github.com/nextbug/shortenerURL
COPY . .

# Build Go Binary
RUN set -ex; \
    CGO_ENABLED=0 GOOS=linux go build -o ./main ./cmd/shortener/main.go;

# Final image, no source code
FROM alpine:latest

# Install Root Ceritifcates
RUN set -ex; \
    apk update; \
    apk add --no-cache \
     ca-certificates

WORKDIR /opt/
COPY --from=src /go/src/github.com/nextbug/shortenerURL/main .

# Run Go Binary
CMD /opt/main