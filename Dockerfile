FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /usdt-app cmd/app/main.go

FROM alpine:3.18

RUN apk update && apk upgrade && apk add --no-cache ca-certificates

COPY --from=builder /usdt-app /usdt-app
COPY migrations /migrations
COPY config.yml /app/config.yml
WORKDIR /app

CMD ["/usdt-app"]
