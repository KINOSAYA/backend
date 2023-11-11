# base go image
FROM golang:1.20.5-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

# CGO=0 means that we do not use C libraries
RUN CGO_ENABLED=0 go build -o extApiservice ./cmd/api

# give extApiservice an executable flag
RUN chmod +x /app/extApiservice

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/extApiservice /app

CMD [ "/app/extApiservice" ]