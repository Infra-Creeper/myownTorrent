package TorrentNet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Payload structure for posting key-value
type PostPayload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Posts a value under a string key to the server
func PostKeyValue(serverURL string, key string, value string) error {
	payload := PostPayload{Key: key, Value: value}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(serverURL+"/put", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}
	return nil
}

func GetValues(serverURL string, key string) ([]string, error) {
	url := fmt.Sprintf("%s/get?key=%s", serverURL, key)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var values []string
	if err := json.Unmarshal(body, &values); err != nil {
		return nil, err
	}
	return values, nil
}
