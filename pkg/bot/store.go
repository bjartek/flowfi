package bot

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// interface to store/load progress
type SubscriptionStore interface {
	// LoadSubscriptions loads the subscriptions from persistent storage.
	LoadSubscriptions() (map[string]SubscriptionData, error)

	// SaveSubscriptions saves the current subscriptions to persistent storage.
	SaveSubscriptions(subscriptions map[string]SubscriptionData) error
}

func (flowFi *FlowFi) StoreProgress() error {
	return flowFi.Store.SaveSubscriptions(flowFi.Subscriptions.pairs)
}

type FileSubscriptionStore struct {
	filename string
	mu       sync.RWMutex
}

func NewFileSubscriptionStore(filename string) *FileSubscriptionStore {
	return &FileSubscriptionStore{filename: filename}
}

// LoadSubscriptions loads subscriptions from the filesystem as pointers.
func (f *FileSubscriptionStore) LoadSubscriptions() (map[string]SubscriptionData, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	file, err := os.Open(f.filename)
	if err != nil {
		if os.IsNotExist(err) {
			// If file does not exist, return an empty map
			return make(map[string]SubscriptionData), nil
		}
		return nil, fmt.Errorf("failed to open subscriptions file: %w", err)
	}
	defer file.Close()

	var subscriptions map[string]SubscriptionData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&subscriptions); err != nil {
		return nil, fmt.Errorf("failed to decode subscriptions: %w", err)
	}

	return subscriptions, nil
}

// SaveSubscriptions saves subscriptions to the filesystem.
func (f *FileSubscriptionStore) SaveSubscriptions(subscriptions map[string]SubscriptionData) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	file, err := os.Create(f.filename)
	if err != nil {
		return fmt.Errorf("failed to create subscriptions file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(subscriptions); err != nil {
		return fmt.Errorf("failed to encode subscriptions: %w", err)
	}

	return nil
}
