package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type KVStore struct {
	mu    sync.RWMutex
	store map[string][]string
}

func NewKVStore() *KVStore {
	return &KVStore{
		store: make(map[string][]string),
	}
}

// Add a value to a key
func (kv *KVStore) putHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	if payload.Key == "" || payload.Value == "" {
		http.Error(w, "key and value required", http.StatusBadRequest)
		return
	}
	kv.mu.Lock()
	kv.store[payload.Key] = append(kv.store[payload.Key], payload.Value)
	kv.mu.Unlock()
	w.WriteHeader(http.StatusOK)
}

// Get all values for a key
func (kv *KVStore) getHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "key required", http.StatusBadRequest)
		return
	}
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	values := kv.store[key]
	json.NewEncoder(w).Encode(values)
}

func main() {
	kv := NewKVStore()
	http.HandleFunc("/put", kv.putHandler)
	http.HandleFunc("/get", kv.getHandler)
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
