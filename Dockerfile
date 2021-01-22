FROM golang:1.14-stretch as build

COPY . /src/tpsloader
WORKDIR /src/tpsloader
RUN go build -o /bin/tpsloader ./internal/main

FROM ubuntu:16.04

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

COPY --from=build /bin/tpsloader /app/
COPY ./config.yml ./config.yml

ENTRYPOINT ["/app/tpsloader"]