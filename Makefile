run: build
	@./bin/goredis --listenAddr :6969
build:
	@go build -o bin/goredis .