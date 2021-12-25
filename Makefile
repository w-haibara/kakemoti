kuirejo: *.go */*.go */*/*.go go.mod
	go mod tidy
	go fmt ./...
	go vet ./...
	gosec ./...
	go build -o kuirejo

.PHONY: test
test: kuirejo
	go test ./...

.PHONY: build-workflow-gen
build-workflow-gen:
	cd workflow && yarn install && tsc index.ts

asl = ""
input = ""
.PHONY: run
run: kuirejo
	node ./workflow/index.js ${asl} > workflow.json
	./kuirejo start-execution \
	--asl workflow.json \
	--input ${input}

