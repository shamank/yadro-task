.PHONY:
.SILENT:

build:
	go build -o ./.bin/app.exe ./cmd/main.go

run: build
	./.bin/app.exe ./tests/case_1.txt

docker-build:
	docker build -t rmaldybaev/yadro --build-arg file_path=./tests/case_1.txt .

docker-run:
	docker run rmaldybaev/yadro