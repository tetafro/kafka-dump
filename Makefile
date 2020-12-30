dep:
	go mod tidy && go mod vendor

test:
	go test ./...

lint:
	golangci-lint run --fix

build:
	go build -mod=vendor -o bin/kafka-dump .

run:
	./bin/kafka-dump
