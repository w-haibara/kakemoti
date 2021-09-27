karage: *.go */*.go go.mod
	go mod tidy
	go fmt ./...
	go vet ./...
	gosec ./...
	go build -o karage

.PHONY: run
run: karage
	./karage start-execution \
	--asl  "./workflow/HelloWorld/statemachine.asl.json" \
	--input "./workflow/HelloWorld/input1.json"