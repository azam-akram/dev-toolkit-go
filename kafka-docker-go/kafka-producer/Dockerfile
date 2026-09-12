FROM golang:1.25 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/producer .

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/producer /producer
EXPOSE 4444
ENTRYPOINT ["/producer"]
