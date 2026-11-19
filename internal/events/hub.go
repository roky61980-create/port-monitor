package events

import (
	"encoding/json"
	"sync"
)

type Hub struct {
	mu sync.RWMutex
	subs map[chan []byte]struct{}
}

func NewHub() *Hub { return &Hub{subs: make(map[chan []byte]struct{})} }

func (h *Hub) Publish(v any) {
	b, err := json.Marshal(v)
	if err != nil { return }
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.subs {
		select { case ch <- b: default: }
	}
}

func (h *Hub) Subscribe() (chan []byte, func()) {
	ch := make(chan []byte, 32)
	h.mu.Lock()
	h.subs[ch] = struct{}{}
	h.mu.Unlock()
	var once sync.Once
	return ch, func() { once.Do(func() { h.mu.Lock(); delete(h.subs, ch); close(ch); h.mu.Unlock() }) }
}
