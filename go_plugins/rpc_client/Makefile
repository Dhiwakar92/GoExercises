binName ?= ./rpcclient

build:
	echo "Building...."
	go mod tidy && env GOOS=darwin GOARCH=arm64 go build -o ${binName} .