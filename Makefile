.PHONY: start build

all: srun

build:
	swag init
	GOOS=linux GOARCH=amd64 go build ./main.go
run: 
	go run main.go
srun:
	swag init
	go run main.go
clean:
	go clean

