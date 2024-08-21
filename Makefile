build:
	go build -C ./cli -o ../bin/kafkito main.go
	go build -C ./server -o ../bin/kafkitoserver main.go
	cp ./config.json ./bin/config.json
