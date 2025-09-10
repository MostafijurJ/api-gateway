package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Item struct {
	SKU   string `json:"sku"`
	Stock int    `json:"stock"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.HandleFunc("/inventory", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode([]Item{{SKU: "ABC-123", Stock: 5}, {SKU: "XYZ-999", Stock: 0}})
	})
	addr := ":9002"
	log.Printf("inventorysvc listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
