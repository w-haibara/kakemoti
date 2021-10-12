karage: *.go */*.go go.mod
	go mod tidy
	go fmt ./...
	go vet ./...
	gosec ./...
	go build -o karage

.PHONY: run
run: run1

.PHONY: run1
run1: karage
	./karage start-execution \
	--asl  "./workflow/HelloWorld/statemachine.asl.json" \
	--input "./workflow/HelloWorld/input1.json"

.PHONY: run2
run2: karage
	./karage start-execution \
	--asl  "./workflow/HelloWorld2/statemachine.asl.json" \
	--input "./workflow/HelloWorld2/input1.json"
