package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Get(key Key) (interface{}, bool)
	Set(key Key, value interface{}) bool
	Clear()
}

type lruCache struct {
	capacity     int
	m            sync.Mutex
	queue        List
	items        map[Key]*ListItem
	reverseItems map[*ListItem]Key
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.m.Lock()
	defer cache.m.Unlock()

	if val, ok := cache.items[key]; ok {
		val.Value = value
		cache.queue.MoveToFront(val)
		return true
	}

	newItem := cache.queue.PushFront(value)
	newItem.Value = value
	cache.items[key] = newItem
	cache.reverseItems[newItem] = key

	if cache.queue.Len() > cache.capacity {
		last := cache.queue.Back()
		delete(cache.items, cache.reverseItems[last])
		delete(cache.reverseItems, last)
		cache.queue.Remove(last)
	}
	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.m.Lock()
	defer cache.m.Unlock()

	if value, ok := cache.items[key]; ok {
		cache.queue.MoveToFront(value)
		return value.Value, true
	}

	return nil, false
}

func (cache *lruCache) Clear() {
	cache.m.Lock()
	defer cache.m.Unlock()
	cache.items = make(map[Key]*ListItem)
	cache.reverseItems = make(map[*ListItem]Key)
	cache.queue = NewList()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity:     capacity,
		queue:        NewList(),
		items:        make(map[Key]*ListItem, capacity),
		reverseItems: make(map[*ListItem]Key),
	}
}
