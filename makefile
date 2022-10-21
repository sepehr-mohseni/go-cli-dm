dev:
	make build
	make run

build:
	go build -o ./bin/main ./main.go

run:
	./bin/main run
