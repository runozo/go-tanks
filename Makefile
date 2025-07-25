.PHONY: wasm

wasm:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w -v" -o ./docs/go-tanks.wasm github.com/runozo/go-tanks

.PHONY: profile

profile:
	go run main.go -cpuprofile output.prof

.PHONY: build

build:
	go build -ldflags="-s -w -v" -o ./cmd/go-tanks

.PHONY: run

run:
	go run main.go