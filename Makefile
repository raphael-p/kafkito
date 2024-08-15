CLI_BINARY_NAME=kafkito
SERVER_BINARY_NAME=kafkitoserver

build:
	go build -C ./cli -o ../bin/${CLI_BINARY_NAME} main.go
	go build -C ./server -o ../bin/${SERVER_BINARY_NAME} main.go
