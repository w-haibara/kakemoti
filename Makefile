kuirejo: *.go */*.go go.mod
	go mod tidy
	go fmt ./...
	go vet ./...
	gosec ./...
	go build -o kuirejo

.PHONY: test
test: kuirejo
	go test ./...
	regresh check
