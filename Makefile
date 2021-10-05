karage: *.go */*.go go.mod
	go mod tidy
	go fmt ./...
	go vet ./...
	gosec ./...
	go build -o karage

.PHONY: run
run: karage
	./karage start-execution \
	--asl  "./workflow/HelloWorld2/statemachine.asl.json" \
	--input "./workflow/HelloWorld2/input1.json"