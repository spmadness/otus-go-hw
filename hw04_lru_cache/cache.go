package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type ListItemValue struct {
	Key   Key
	Value interface{}
}

var mu sync.Mutex

func (l *lruCache) Set(key Key, value interface{}) bool {
	mu.Lock()
	defer mu.Unlock()

	val, ok := l.items[key]
	if ok {
		val.Value.(*ListItemValue).Value = value
		l.queue.MoveToFront(val)
		return true
	}

	l.items[key] = l.queue.PushFront(&ListItemValue{
		Key:   key,
		Value: value,
	})

	if l.queue.Len() > l.capacity {
		backKey := l.queue.Back().Value.(*ListItemValue).Key
		delete(l.items, backKey)

		l.queue.Remove(l.queue.Back())
	}

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	mu.Lock()
	defer mu.Unlock()

	if val, ok := l.items[key]; ok {
		l.queue.MoveToFront(val)
		return val.Value.(*ListItemValue).Value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	mu.Lock()
	defer mu.Unlock()

	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
