build: format
	go build

format:
	gofmt -w server.go
