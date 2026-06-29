# Build stage
FROM golang:1.21-alpine3.19 AS builder

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
ENV CGO_ENABLED=1

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -v -ldflags="-s -w" -o server ./cmd/server

# Final stage
FROM alpine:3.19

RUN apk add --no-cache tzdata ca-certificates

# Set timezone to Asia/Shanghai
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server /app/server
COPY --from=builder /app/.env /app/.env

# Create logs directory
RUN mkdir -p /app/logs

EXPOSE 8888

CMD ["/app/server"]
