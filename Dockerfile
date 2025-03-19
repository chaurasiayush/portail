# Dockerfile
FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o portail ./cmd/portail

# Final image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/portail .
COPY config.yaml .

EXPOSE 8080
EXPOSE 5353/udp

ENTRYPOINT ["./portail"]
