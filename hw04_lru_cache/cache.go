package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
}

type lruCache struct {
	mutex    sync.Mutex
	items    map[Key]*ListItem
	queue    List
	capacity int
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem),
	}
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	cacheValue := cacheItem{key: key, value: value}
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	if item, ok := lc.items[key]; ok {
		item.Value = cacheValue
		lc.queue.MoveToFront(item)
		return true
	}
	if lc.queue.Len() >= lc.capacity {
		back := lc.queue.Back()
		value := back.Value.(cacheItem)
		delete(lc.items, value.key)
		lc.queue.Remove(back)
	}
	item := lc.queue.PushFront(cacheValue)
	lc.items[key] = item
	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	if item, ok := lc.items[key]; ok {
		lc.queue.MoveToFront(item)
		cacheValue := item.Value.(cacheItem)
		return cacheValue.value, true
	}
	return nil, false
}

func (lc *lruCache) Clear() {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()
	for key, item := range lc.items {
		lc.queue.Remove(item)
		delete(lc.items, key)
	}
}