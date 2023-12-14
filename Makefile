.PHONY: webasm

webasm:
	GOOS=js GOARCH=wasm go build -ldflags="-s -w -v" -o ./docs/go-tanks.wasm github.com/runozo/go-tanks

.PHONY: run

run:
	go run main.go