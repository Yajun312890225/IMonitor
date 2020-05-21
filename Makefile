.PHONY: start build

all: run

build:
	GOOS=linux GOARCH=amd64 go build ./main.go
run: 
	go run main.go
srun:
	swag init
	go run main.go
clean:
	go clean

