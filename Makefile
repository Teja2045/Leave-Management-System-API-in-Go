start:
	go run server.go

hello:
	echo "hello"

build:
	go build -o bin/server server.go

test:
	go test