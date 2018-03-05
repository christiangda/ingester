package server

import (
	"sync"

	"github.com/google/uuid"
)

// ConnectionStore store in memory
// reference: http://dnaeon.github.io/concurrent-maps-and-slices-in-go/
type ConnectionStore struct {
	sync.RWMutex
	items map[uuid.UUID]*Connection
}

// ConnectionStoreItem items
type ConnectionStoreItem struct {
	Key   uuid.UUID
	Value *Connection
}

// NewConnectionStore a new ConnectionStore
func NewConnectionStore() *ConnectionStore {
	cs := &ConnectionStore{
		items: make(map[uuid.UUID]*Connection),
	}

	return cs
}

// Set a key in a Connection storage
func (cs *ConnectionStore) Set(key uuid.UUID, value *Connection) bool {
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

// Get a key from a Connection storage
func (cs *ConnectionStore) Get(key uuid.UUID) (*Connection, bool) {
	cs.RLock()
	defer cs.RUnlock()

	value, ok := cs.items[key]
	return value, ok
}

// Put a new wlwmwnt in store and return the new key
func (cs *ConnectionStore) Put(value *Connection) uuid.UUID {
	key, _ := uuid.NewUUID()

	cs.Set(key, value)
	return key
}

// Delete a key from a Connection storage
func (cs *ConnectionStore) Delete(key uuid.UUID) {
	cs.Lock()
	delete(cs.items, key)
	cs.Unlock()
}

// Count return number of items in Connection storage
func (cs *ConnectionStore) Count() int {
	cs.RLock()
	defer cs.RUnlock()
	return len(cs.items)
}

// Iter iterates over the items in a cConnection storage
// Each item is sent over a channel, so that
// we can iterate over the Connection storage using the builtin range keyword
func (cs *ConnectionStore) Iter() <-chan ConnectionStoreItem {
	c := make(chan ConnectionStoreItem)

	go func() {
		cs.Lock()
		defer cs.Unlock()

		for k, v := range cs.items {
			c <- ConnectionStoreItem{k, v}
		}
		close(c)
	}()

	return c
}
