build:
	@go build -o ./bin/go-chain

run: build
	./bin/go-chain

test:
	go test ./...