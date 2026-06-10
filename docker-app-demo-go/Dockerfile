FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o app .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
