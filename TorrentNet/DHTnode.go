package TorrentNet

import (
	"encoding/hex"
	"errors"
	"sync"

	"github.com/purehyperbole/dht"

	"time"
)

func GetSeeds(node *dht.DHT, hash string, timeout time.Duration) ([]string, error) {
	var seeds []string
	var mu sync.Mutex
	var lookupErr error
	done := make(chan bool)
	key, decodeErr := hex.DecodeString(hash)
	if decodeErr != nil {
		return []string{}, decodeErr
	}
	go func() {
		node.Find(key, func(value []byte, err error) {
			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				if lookupErr == nil {
					lookupErr = err
				}
				return
			}

			// Copy value before storing (DHT reuses buffers)
			copied := make([]byte, len(value))
			copy(copied, value)
			ip := string(copied)
			seeds = append(seeds, ip)
		})

		// Wait for lookup to complete
		time.Sleep(timeout)
		done <- true
	}()

	<-done
	return seeds, lookupErr
}

func PostSeed(d *dht.DHT, hash string, ipAddr string, ttl time.Duration) error {
	key, decodeErr := hex.DecodeString(hash)
	if decodeErr != nil {
		return decodeErr
	}
	value := []byte(ipAddr)
	var funcErr error = nil

	d.Store(key, value, ttl, func(err error) {
		if err != nil {
			funcErr = errors.New("ERROR: Unable to store " + hash + "=" + ipAddr + "on DHT")
		}
	})
	return funcErr
}
