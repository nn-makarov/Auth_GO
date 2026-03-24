FROM golang:1.26-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o app .

FROM alpine:latest

RUN apk add --no-cache sqlite

WORKDIR /root/

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]