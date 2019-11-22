dep:
	go mod tidy && go mod vendor

build:
	go build -mod=vendor -o bin/kafka-dump .

run:
	./bin/kafka-dump
