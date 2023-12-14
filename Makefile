.PHONY: webasm

webasm:
	GOOS=js GOARCH=wasm go build -o ./out/go-tanks.wasm github.com/runozo/go-tanks

.PHONY: run

run:
	go run main.go