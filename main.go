package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// maximum size of the KVS
var CAPACITY = 1000

// thread safe and goroutine safe are assumed to be interchangable here.
// https://groups.google.com/g/golang-nuts/c/ZA0IK1k6UVk
// Theoretically, goroutine safe concurrency === thread safe concurrency

type KVS[V any] []*Node[V]

type ConcurrentKVS[V any] struct {
	kvs KVS[V]
	// https://pkg.go.dev/sync#RWMutex
	// Can be held by arbitrary # of readers, single writer.
	Mu sync.RWMutex
	// for edge-triggered persistence: persist if conditional evaluates to true
	cond func(kvs *KVS[V]) bool
}

type Node[V any] struct {
	key   string
	value V
}

// Instantiate new concurrent map
func New[V any](f func(*KVS[V]) bool) *ConcurrentKVS[V] {
	m := make([]*Node[V], CAPACITY)
	return &ConcurrentKVS[V]{kvs: m, cond: f}
}

func (m *ConcurrentKVS[V]) Get(key string) V {
	m.Mu.RLock()
	defer m.Mu.RUnlock()

	index := hashFunction(key)
	return m.kvs[index].value
}

func (m *ConcurrentKVS[V]) Put(key string, value V) error {
	m.Mu.Lock()
	defer m.Mu.Unlock()

	index := hashFunction(key)
	// search for collision first
	// fix: just fails, no resolution path here?
	if m.kvs[index] != nil {
		return fmt.Errorf("collision at index %d when trying to put key %s", index, key)
	}
	m.kvs[index] = &Node[V]{key, value}

	// Persist to disk on every put if conditional is met

	if m.cond(&m.kvs) {
		m.Persist()
	}

	return nil
}

func (m *ConcurrentKVS[V]) Persist() error {
	b, err := json.MarshalIndent(m.kvs, " ", " ")
	if err != nil {
		return fmt.Errorf("%v: could not marshal key-value store to json", err)
	}
	os.WriteFile("kvs.json", b, 0o644)
	return nil
}

// Collisions are a potential issue with modular hashing.
// Can use linked lists or a better hashing function?
func hashFunction(in string) int {
	var ret int
	for i := 0; i < len(in); i++ {
		// 31 is what java uses :)
		ret = (31*ret + int(in[i])) % CAPACITY
	}
	return ret
}
