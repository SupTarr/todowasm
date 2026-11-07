# TodoMVC written in Go and WebAssembly

This project is a small TodoMVC-style application with:

- A Go HTTP server for the frontend and todo API.
- A Go WebAssembly frontend compiled to `wasm/assets/lib.wasm`.
- An in-memory todo store, so data is intentionally reset whenever the server
  restarts.
- A login screen that retrieves a local development token before loading todos.

## Requirements

- Go 1.24 or newer
- A browser with WebAssembly support

No JavaScript package manager or frontend build tool is required. The browser
runtime used by Go is already checked into `wasm/assets/wasm_exec.js`.

## Build the WebAssembly frontend

From the repository root, run:

```sh
GOOS=js GOARCH=wasm go build -o wasm/assets/lib.wasm ./wasm
```

Or use the Makefile:

```sh
make wasm
```

Run this again whenever the Go code under `wasm/` or `todo/` changes.

## Run locally

Build the WebAssembly frontend, then start the server:

```sh
make wasm
go run .
```

Open <http://localhost:8080> in a browser. The server hosts both the static
frontend and the API from the same origin.

The server accepts these options:

```text
-listen string   address to listen on (default ":8080")
-dir string      directory containing the frontend assets (default "./wasm/assets")
```

For example:

```sh
go run . -listen :9090 -dir ./wasm/assets
```

The Makefile also provides a shorthand:

```sh
make run
```

## API

The API is served by the same process as the frontend:

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/tokenz` | Returns the local development token. |
| `GET` | `/todos` | Lists all todos. |
| `POST` | `/todos` | Creates a todo from a JSON body such as `{"text":"Buy milk"}`. |
| `DELETE` | `/todos/:id` | Deletes a todo by numeric ID. |

Todo responses use this shape:

```json
{
  "id": 1,
  "text": "Buy milk"
}
```

The todo endpoints require the authorization header used by the WebAssembly
client. This authentication flow is intended for local development and is not
production security.

## Project layout

```text
.
├── main.go              # HTTP server, static files, and API routes
├── todo/
│   ├── api.go           # Browser-side API calls
│   └── todo.go          # Browser UI and todo interactions
├── wasm/
│   ├── main.go          # WebAssembly entry point
│   └── assets/          # HTML, CSS, Go runtime, and compiled WASM
├── Makefile              # Build and run shortcuts
└── go.mod                # Go module and dependencies
```

## Development notes

- Todos are stored only in memory; there is no database or persistence layer.
- The generated `wasm/assets/lib.wasm` file must exist before opening the
  frontend.
- The server serves files from `wasm/assets` by default, so run commands from
  the repository root or pass an explicit `-dir` path.
