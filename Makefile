.PHONY: start build

all: run

build:
	swag init
	GOOS=linux GOARCH=amd64 go build ./main.go && mv main start
run:
	swag init
	go run main.go
clean:
	go clean
swagger:
	swag init

