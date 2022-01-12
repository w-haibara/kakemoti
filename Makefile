kakemoti: *.go */*.go */*/*.go go.mod
	go mod tidy
	go fmt ./...
	go vet ./...
	gosec -exclude-dir=_workflow ./...
	go build -o kakemoti

_workflow/index.js: _workflow/*.ts
	cd _workflow && yarn install && tsc

asl = ""
.PHONY: asl-gen
asl-gen: _workflow/index.js
	node ./_workflow/index.js ${asl} > _workflow/asl/${asl}.asl.json

.PHONY: asl-gen-all
asl-gen-all: _workflow/index.js
	mkdir -p ./_workflow/asl
	node ./_workflow/index.js list | while read -r a; do eval "make asl-gen asl=$$a"; done

.PHONY: clean
clean:
	rm ./_workflow/asl/*

input = ""
.PHONY: run
run: kakemoti
	./kakemoti start-execution \
		--asl _workflow/asl/${asl}.asl.json \
		--input ${input}

.PHONY: test
test: kakemoti asl-gen-all
	go test -count=1 ./...
