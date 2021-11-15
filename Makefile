kuirejo: *.go */*.go go.mod
	go mod tidy
	go fmt ./...
	go vet ./...
	gosec ./...
	go build -o kuirejo

.PHONY: run
run: run5

.PHONY: run1
run1: kuirejo
	./kuirejo start-execution \
	--asl  "./workflows/HelloWorld/statemachine.asl.json" \
	--input "./workflows/HelloWorld/input1.json"

.PHONY: run2
run2: kuirejo
	./kuirejo start-execution \
	--asl  "./workflows/HelloWorld2/statemachine.asl.json" \
	--input "./workflows/HelloWorld2/input1.json"

.PHONY: run3
run3: kuirejo
	./kuirejo start-execution \
	--asl  "./workflows/task-script1/statemachine.asl.json" \
	--input "./workflows/task-script1/input1.json"

.PHONY: run4
run4: kuirejo
	./kuirejo start-execution \
	--asl  "./workflows/task-script2/statemachine.asl.json" \
	--input "./workflows/task-script2/input1.json"

.PHONY: run5
run5: kuirejo
	./kuirejo start-execution \
	--asl  "./workflows/task-script3/statemachine.asl.json" \
	--input "./workflows/task-script3/input1.json"
