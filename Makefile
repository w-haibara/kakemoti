kakemoti: *.go */*.go */*/*.go go.mod
	go mod tidy
	go fmt ./...
	go vet ./...
	gosec -exclude-dir=_workflow ./...
	go build -o kakemoti

_workflow/index.js: _workflow/*.ts
	cd _workflow && yarn install && tsc

_workflow/asl/%.asl.json: _workflow/index.js
	node ./_workflow/index.js $@

.PHONY: asl-gen-all
asl-gen-all: _workflow/index.js
	mkdir -p ./_workflow/asl
	node ./_workflow/index.js list | while read -r a; do eval "make _workflow/asl/$$a.asl.json"; done

.PHONY: test
test: asl-gen-all
	go test -count=1 ./...

input = ""
.PHONY: run
run: kakemoti
	./kakemoti start-execution \
		--asl _workflow/asl/${asl}.asl.json \
		--input ${input}

.PHONY: clean
clean:
	rm ./_workflow/asl/*
