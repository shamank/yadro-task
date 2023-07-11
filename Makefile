.PHONY:
.SILENT:

build:
	go build -o ./.bin/app.exe ./cmd/main.go

run: build
	./.bin/app.exe $(FILE)

docker-build:
	docker build -t rmaldybaev/yadro .

docker-run:
	docker run rmaldybaev/yadro $(FILE)

docker: docker-build docker-run