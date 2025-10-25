package TorrentNet

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/purehyperbole/dht"
)

// CustomStorage tracks all stored keys
type CustomStorage struct {
	data map[string]*StoredValue
	mu   sync.RWMutex
}

// An Data Entry unit of an DHT
type StoredValue struct {
	key     []byte
	value   []byte
	created time.Time
	ttl     time.Duration
}

// Creates New Custom Storage unit
func NewCustomStorage() *CustomStorage {
	return &CustomStorage{
		data: make(map[string]*StoredValue),
	}
}

// Sets a value on the CustomeStorage
func (s *CustomStorage) Set(key, value []byte, created time.Time, ttl time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[string(key)] = &StoredValue{
		key:     key,
		value:   value,
		created: created,
		ttl:     ttl,
	}

	log.Printf("📥 Received and stored: %s -> %s", string(key), string(value))
	return true
}

// Gets the value of the custome storage
func (s *CustomStorage) Get(key []byte, from time.Time) ([]*dht.Value, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if val, ok := s.data[string(key)]; ok {
		if time.Since(val.created) < val.ttl {
			return []*dht.Value{
				{
					Key:     val.key,
					Value:   val.value,
					Created: val.created,
				},
			}, true
		}
	}
	return nil, false
}

// Iterates over the customStorage and executes cb till it returns false
func (s *CustomStorage) Iterate(cb func(value *dht.Value) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, val := range s.data {
		if time.Since(val.created) < val.ttl {
			if !cb(&dht.Value{
				Key:     val.key,
				Value:   val.value,
				Created: val.created,
			}) {
				return
			}
		}
	}
}

// Prints all key-value pair of custom storage
func (s *CustomStorage) PrintAll() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	fmt.Println("\n📋 Locally stored key-value pairs:")

	for _, val := range s.data {
		if time.Since(val.created) < val.ttl {
			fmt.Printf("   %s = %s\n", string(val.key), string(val.value))
			count++
		}
	}

	if count == 0 {
		fmt.Println("   (No keys stored locally yet)")
	}
	fmt.Println()
}
