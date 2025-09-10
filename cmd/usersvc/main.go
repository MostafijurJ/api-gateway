package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]User{{ID: 1, Name: "Alice"}, {ID: 2, Name: "Bob"}})
	})
	addr := ":9001"
	log.Printf("usersvc listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
