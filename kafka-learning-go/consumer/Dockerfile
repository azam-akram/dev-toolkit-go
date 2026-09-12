FROM golang:1.25 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOFLAGS=-mod=vendor go build -o /out/consumer .

FROM gcr.io/distroless/static-debian12
COPY --from=build /out/consumer /consumer
EXPOSE 5555
ENTRYPOINT ["/consumer"]
