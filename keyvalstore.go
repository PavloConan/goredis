package main

import "sync"

type KeyValStore struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func NewKeyValStore() *KeyValStore {
	return &KeyValStore{
		data: make(map[string][]byte),
	}
}

func (kv *KeyValStore) Set(key, val []byte) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	kv.data[string(key)] = val
	return nil
}

func (kv *KeyValStore) Get(key []byte) ([]byte, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	val, ok := kv.data[string(key)]

	return val, ok
}
