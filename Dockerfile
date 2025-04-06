FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server .



FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/server /app/server

EXPOSE 9000

CMD ["/app/server"]
