package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Item struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Store struct {
	mu    sync.Mutex
	next  int
	items map[int]Item
}

func NewStore() *Store {
	return &Store{next: 1, items: map[int]Item{}}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func NewMux(s *Store) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /items", func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		list := make([]Item, 0, len(s.items))
		for i := 1; i < s.next; i++ {
			if it, ok := s.items[i]; ok {
				list = append(list, it)
			}
		}
		writeJSON(w, 200, list)
	})

	mux.HandleFunc("POST /items", func(w http.ResponseWriter, r *http.Request) {
		var in Item
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Name == "" {
			writeJSON(w, 400, map[string]string{"error": "name wajib diisi"})
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		in.ID = s.next
		s.next++
		s.items[in.ID] = in
		writeJSON(w, 201, in)
	})

	mux.HandleFunc("PUT /items/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		var in Item
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.Name == "" {
			writeJSON(w, 400, map[string]string{"error": "name wajib diisi"})
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		if _, ok := s.items[id]; !ok {
			writeJSON(w, 404, map[string]string{"error": "not found"})
			return
		}
		in.ID = id
		s.items[id] = in
		writeJSON(w, 200, in)
	})

	mux.HandleFunc("DELETE /items/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		s.mu.Lock()
		defer s.mu.Unlock()
		if _, ok := s.items[id]; !ok {
			writeJSON(w, 404, map[string]string{"error": "not found"})
			return
		}
		delete(s.items, id)
		writeJSON(w, 200, map[string]int{"deleted": id})
	})

	return mux
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, NewMux(NewStore())))
}