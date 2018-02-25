package server

import (
	"sync"

	"github.com/google/uuid"
)

// ClientStore store in memory
// reference: http://dnaeon.github.io/concurrent-maps-and-slices-in-go/
type ClientStore struct {
	sync.RWMutex
	items map[uuid.UUID]*Client
}

// ClientStoreItem items
type ClientStoreItem struct {
	Key   uuid.UUID
	Value *Client
}

// NewClientStore a new ClientStore
func NewClientStore() *ClientStore {
	cs := &ClientStore{
		items: make(map[uuid.UUID]*Client),
	}

	return cs
}

// Set a key in a client storage
func (cs *ClientStore) Set(key uuid.UUID, value *Client) bool {
	cs.Lock()

	_, present := cs.items[key]
	if present {
		cs.Unlock()
		return false
	}

	cs.items[key] = value
	cs.Unlock()
	return true
}

// Get a key from a client storage
func (cs *ClientStore) Get(key uuid.UUID) (*Client, bool) {
	cs.RLock()
	defer cs.RUnlock()

	value, ok := cs.items[key]
	return value, ok
}

// Put a new wlwmwnt in store and return the new key
func (cs *ClientStore) Put(value *Client) uuid.UUID {
	key, _ := uuid.NewUUID()

	cs.Set(key, value)
	return key
}

// Delete a key from a client storage
func (cs *ClientStore) Delete(key uuid.UUID) {
	cs.Lock()
	delete(cs.items, key)
	cs.Unlock()
}

// Count return number of items in client storage
func (cs *ClientStore) Count() int {
	cs.RLock()
	defer cs.RUnlock()
	return len(cs.items)
}

// Iter iterates over the items in a cclient storage
// Each item is sent over a channel, so that
// we can iterate over the client storage using the builtin range keyword
func (cs *ClientStore) Iter() <-chan ClientStoreItem {
	c := make(chan ClientStoreItem)

	go func() {
		cs.Lock()
		defer cs.Unlock()

		for k, v := range cs.items {
			c <- ClientStoreItem{k, v}
		}
		close(c)
	}()

	return c
}
