package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	listen = flag.String("listen", ":8080", "listen address")
	dir    = flag.String("dir", "./wasm/assets", "directory to serve")
)

type serverTodo struct {
	ID        uint   `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

type todoStore struct {
	mu     sync.RWMutex
	nextID uint
	todos  []serverTodo
}

func main() {
	flag.Parse()
	store := &todoStore{nextID: 1}
	mux := http.NewServeMux()
	mux.HandleFunc("/tokenz", tokenHandler)
	mux.HandleFunc("/todos", store.todosHandler)
	mux.HandleFunc("/todos/", store.todoHandler)
	mux.Handle("/", http.FileServer(http.Dir(*dir)))

	log.Printf("listening on %q...", *listen)
	log.Fatal(http.ListenAndServe(*listen, mux))
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": "local-dev-token"})
}

func (s *todoStore) todosHandler(w http.ResponseWriter, r *http.Request) {
	if !authorized(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.mu.RLock()
		todos := append([]serverTodo(nil), s.todos...)
		s.mu.RUnlock()
		writeJSON(w, http.StatusOK, todos)
	case http.MethodPost:
		var input serverTodo
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil || strings.TrimSpace(input.Text) == "" {
			http.Error(w, "text is required", http.StatusBadRequest)
			return
		}

		s.mu.Lock()
		input.ID = s.nextID
		s.nextID++
		input.Text = strings.TrimSpace(input.Text)
		s.todos = append(s.todos, input)
		s.mu.Unlock()
		writeJSON(w, http.StatusCreated, input)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *todoStore) todoHandler(w http.ResponseWriter, r *http.Request) {
	if !authorized(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method != http.MethodDelete && r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseUint(strings.TrimPrefix(r.URL.Path, "/todos/"), 10, 0)
	if err != nil || id == 0 {
		http.Error(w, "invalid todo id", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for i, todo := range s.todos {
		if todo.ID == uint(id) {
			if r.Method == http.MethodPatch {
				var update struct {
					Completed *bool `json:"completed"`
				}
				if err := json.NewDecoder(r.Body).Decode(&update); err != nil || update.Completed == nil {
					http.Error(w, "completed is required", http.StatusBadRequest)
					return
				}
				todo.Completed = *update.Completed
				s.todos[i] = todo
				writeJSON(w, http.StatusOK, todo)
				return
			}
			s.todos = append(s.todos[:i], s.todos[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "todo not found", http.StatusNotFound)
}

func authorized(r *http.Request) bool {
	return r.Header.Get("Authorization") == "Bearer local-dev-token"
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
