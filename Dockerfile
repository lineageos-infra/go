FROM golang:1.15 as builder

COPY . /app
WORKDIR /app

RUN go build -o /usr/local/bin/lineage

FROM ubuntu:20.04

COPY --from=builder /usr/local/bin/lineage /usr/local/bin/lineage