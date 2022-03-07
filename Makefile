## Invoke linter to promote Go best practices
lint:
	golangci-lint run ./...

# Run go fmt against code
fmt:
	go fmt ./...

## Inspects the source code for suspicious constructs
vet:
	go vet ./...

# Build tasty binary
build: fmt vet lint
	go build -o bin/tasty main.go

test: fmt vet lint
	go test ./... --cover

release: fmt vet lint
	env GOOS=darwin GOARCH=amd64 go build -o bin/tasty-darwin-amd64
	env GOOS=darwin GOARCH=arm64 go build -o bin/tasty-darwin-arm64
	env GOOS=linux GOARCH=amd64 go build -o bin/tasty-linux-amd64 ; go build -o bin/tasty-linux-x86_64
	env GOOS=linux GOARCH=arm64 go build -o bin/tasty-linux-arm64 ; go build -o bin/tasty-linux-aarch64
