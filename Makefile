kakemoti: *.go */*.go */*/*.go go.mod
	go mod tidy
	go fmt ./...
	go vet ./...
	gosec -exclude-dir=_workflow ./...
	go build -o kakemoti

.PHONY: test
test: kakemoti build-workflow-gen
	go test -count=1 ./...

.PHONY: build-workflow-gen
build-workflow-gen: _workflow/*
	cd _workflow && yarn install && tsc

asl = ""
.PHONY: workflow-gen
workflow-gen: kakemoti
	node ./_workflow/index.js ${asl} > workflow.json

input = ""
.PHONY: workflow-run
workflow-run: kakemoti workflow-gen
	./kakemoti start-execution \
		--asl workflow.json \
		--input ${input}
