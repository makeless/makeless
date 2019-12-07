lint:
	- golangci-lint run

go-build:
	GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o serve-darwin
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o serve-linux