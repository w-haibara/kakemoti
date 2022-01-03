kuirejo: *.go */*.go */*/*.go go.mod
	go mod tidy
	go fmt ./...
	go vet ./...
	gosec -exclude-dir=_workflow ./...
	go build -o kuirejo

.PHONY: test
test: kuirejo
	go test -count=1 ./...

.PHONY: build-workflow-gen
build-workflow-gen:
	cd _workflow && yarn install && tsc

asl = ""
.PHONY: workflow-gen
workflow-gen: kuirejo
	node ./_workflow/index.js ${asl} > workflow.json

input = ""
.PHONY: run-workflow
workflow-run: kuirejo workflow-gen
	./kuirejo start-execution \
		--asl workflow.json \
		--input ${input}
