all: build

build:
	go build -o crawler && ./crawler

run:
	go run .

test:
	go test -count=1 -v ./...
