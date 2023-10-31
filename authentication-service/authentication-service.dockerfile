# base go image
FROM golang:1.20.5-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

# CGO=0 means that we do not use C libraries
RUN CGO_ENABLED=0 go build -o authApp ./cmd/api

# give authApp an executable flag
RUN chmod +x /app/authApp

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/authApp /app

CMD [ "/app/authApp" ]