FROM golang:latest

ARG file_path

WORKDIR /app

COPY . .

RUN go build -o ./.bin/app.exe ./cmd/main.go

CMD ["./.bin/app.exe", "$file_path"]