package main

import (
	"sync"
)

type Broadcaster[T any] struct {
	mu        sync.RWMutex
	channel   chan T
	listeners map[int]chan T
	nextID    int
}

func NewBroadcaster[T any]() *Broadcaster[T] {
	return &Broadcaster[T]{
		channel:   make(chan T),
		listeners: make(map[int]chan T),
		nextID:    0,
	}
}

func (b *Broadcaster[T]) Listen() <-chan T {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan T)
	b.listeners[b.nextID] = ch

	b.nextID += 1

	return ch
}

func (b *Broadcaster[T]) Broadcast(msg T) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, listener := range b.listeners {
		select {
		case listener <- msg:
		default:
		}
	}
}

func (b *Broadcaster[T]) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, listener := range b.listeners {
		close(listener)
	}
	b.listeners = make(map[int]chan T)
}
