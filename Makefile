linter:
	- golangci-lint run

go-build:
	GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o makeless-darwin
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o makeless-linux
	chmod +x makeless-darwin
	chmod +x makeless-linux