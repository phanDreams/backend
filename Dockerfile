FROM golang:1.24.1-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o ./bin/main cmd/main.go

FROM alpine:3.21

RUN adduser -D appuser
USER appuser

WORKDIR /home/appuser
COPY --from=builder /app/bin/main /usr/local/bin/

EXPOSE 8080

CMD ["/usr/local/bin/main"]