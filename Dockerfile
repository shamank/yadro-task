FROM golang:latest as build

WORKDIR /app

COPY . .

RUN go build -o ./.bin/app.exe ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/.bin/app.exe ./.bin/app.exe
COPY ./tests/ ./tests

ENTRYPOINT ["./.bin/app.exe"]
