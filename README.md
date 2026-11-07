# TodoMVC written in Go

## How to build

> GOARCH=wasm GOOS=js go build -o wasm/assets/lib.wasm wasm/main.go

### or use Makefile

> make wasm

## How to run

> go run main.go

Then open http://localhost:8080. The server serves the WebAssembly frontend and
the local todo API used by it. Todos are kept in memory and reset on restart.
