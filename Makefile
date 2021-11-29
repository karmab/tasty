build:
	go build tasty && mv tasty tasty-darwin-amd64
	env GOOS=linux GOARCH=amd64 go build tasty && mv tasty tasty-linux-amd64

test:
	go test ./... --cover
