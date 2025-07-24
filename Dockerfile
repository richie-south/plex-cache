FROM golang:1.24.5 AS builder

WORKDIR /app

COPY . .
COPY .env .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o plex-cache ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/plex-cache .
COPY --from=builder /app/.env .

CMD ["./plex-cache"]
