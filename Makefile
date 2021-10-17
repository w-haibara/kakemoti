karage: *.go */*.go go.mod
	go mod tidy
	go fmt ./...
	go vet ./...
	gosec ./...
	go build -o karage

.PHONY: run
run: run3

.PHONY: run1
run1: karage
	./karage start-execution \
	--asl  "./workflows/HelloWorld/statemachine.asl.json" \
	--input "./workflows/HelloWorld/input1.json"

.PHONY: run2
run2: karage
	./karage start-execution \
	--asl  "./workflows/HelloWorld2/statemachine.asl.json" \
	--input "./workflows/HelloWorld2/input1.json"

.PHONY: run3
run3: karage
	./karage start-execution \
	--asl  "./workflows/task-script1/statemachine.asl.json" \
	--input "./workflows/task-script1/input1.json"
