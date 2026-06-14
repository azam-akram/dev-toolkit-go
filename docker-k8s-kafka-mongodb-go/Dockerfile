FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd

FROM alpine:latest 

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD [ "./server" ]