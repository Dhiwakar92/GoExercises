binName ?= ./add

build:
	echo "Building...."
	go mod tidy && env GOOS=darwin GOARCH=arm64 go build -o ${binName} .
	mv ${binName} /Users/dhiwakarravikumar/experimental/reviewing_go_plugins/rpc_client