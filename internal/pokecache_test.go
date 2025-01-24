package internal

import (
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	cache := NewCache(time.Duration(5) * time.Second)
	cache.Add("google.com", []byte{0, 4, 8})
	data, ok := cache.Get("google.com")
	if !ok {
		t.Error("couldn't get data from cache")
	}
	if len(data) != len([]byte{0, 4, 8}) {
		t.Errorf("data from cache has length %d, original data has length %d", len(data), len([]byte{0, 4, 8}))
	}
	for i := range data {
		if data[i] != []byte{0, 4, 8}[i] {
			t.Error("original value and value in cache are different")
		}
	}
}

func TestReapLoop(t *testing.T) {
	cache := NewCache(time.Duration(5) * time.Second)
	cache.Add("google.com", []byte{0, 4, 8})
	time.Sleep(time.Duration(7) * time.Second)
	_, ok := cache.Get("google.com")
	if ok {
		t.Error("got data from cache that should've been invalid")
	}
}
