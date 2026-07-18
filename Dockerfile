FROM golang:1.26.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app ./cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/app .

COPY --from=builder /app/config ./config/

EXPOSE 8080

ENTRYPOINT ["./app"]