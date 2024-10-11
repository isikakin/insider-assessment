FROM golang:1.23.2-alpine3.20 AS builder
WORKDIR /opt
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /opt/configs/config.yaml /app/
COPY --from=builder /opt/main /app/main
WORKDIR /app
ENTRYPOINT ["/app/main"]