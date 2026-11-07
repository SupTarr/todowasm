.PHONY: wasm run

wasm:
	GOARCH=wasm GOOS=js go build -o wasm/assets/lib.wasm wasm/main.go

run:
	go run main.go