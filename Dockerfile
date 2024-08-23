# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

COPY *.go ./

RUN go build -v -o /prusa_metrics_handler

FROM alpine:latest

COPY --from=builder /prusa_metrics_handler .

EXPOSE 10011

ENTRYPOINT ["/prusa_metrics_handler"]