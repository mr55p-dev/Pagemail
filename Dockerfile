# Build stage
FROM golang:1.23-bookworm AS builder

WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y make curl git

COPY Makefile .
ENV GOBIN /app/bin
RUN make bin/tailwindcss
# RUN make bin/dbmate
RUN make $GOBIN/templ
RUN make $GOBIN/sqlc

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application using Makefile
# The Makefile expects tools to be in specific places or GOBIN.
# We installed templ and sqlc to GOBIN (/go/bin) which is in PATH.
# We downloaded tailwindcss to `./bin/tailwindcss` as expected by Makefile.
RUN make server

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/pagemail /app/pagemail
# COPY --from=builder /app/migrate .

# Copy assets if needed (css is embedded or needed?) 
# Looking at Makefile, css is compiled to assets/css/main.css. 
# We should verify if application serves assets from disk.
COPY --from=builder /app/assets ./assets

# Expose port
EXPOSE 8080

CMD ["ls", "/app"]
