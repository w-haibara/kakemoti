.PHONY: test
test: asl-gen-all
	go fmt ./...
	go vet ./...
	gosec -exclude-dir=_workflow ./...
	go test -v -count=1 ./...

_workflow/index.js: _workflow/*.ts
	cd _workflow && yarn install && tsc

_workflow/asl/%.asl.json: _workflow/index.js
	node ./_workflow/index.js $@

.PHONY: asl-gen-all
asl-gen-all: _workflow/index.js
	mkdir -p ./_workflow/asl
	node ./_workflow/index.js list | while read -r a; do eval "make _workflow/asl/$$a.asl.json"; done

.PHONY: clean
clean:
	rm ./_workflow/asl/*
