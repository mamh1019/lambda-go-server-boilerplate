FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./main.go

FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates && update-ca-certificates

COPY --from=builder /app/server /app/server

EXPOSE 8080

ENV GIN_MODE=release

CMD ["/app/server"]

