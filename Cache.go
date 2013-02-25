package main

import (
  "container/list"
	"fmt"
	"sync"
	"time"
)

// func main() {
// 	cache := NewCache(12)
// 	cache.Put("a", "Hi a")
// 	time.Sleep(1 * time.Second)
// 	cache.Put("b", "Hi b")
// 	time.Sleep(1 * time.Second)
// 	cache.Put("c", "Hi c")

// 	for e := cache.list.Front(); e != nil; e = e.Next() {
// 		info := e.Value.(*ElementInfo)
// 		fmt.Printf("k:%s \t v:%s \t t:%s \t %d \n", info.Key, info.Value.(string), info.TimeAccessed.String(), info.TimeAccessed.Unix())
// 	}

// 	time.Sleep(1 * time.Second)
// 	cache.Get("a")
// 	for e := cache.list.Front(); e != nil; e = e.Next() {
// 		info := e.Value.(*ElementInfo)
// 		fmt.Printf("k:%s \t v:%s \t t:%s \t %d \t \n", info.Key, info.Value.(string), info.TimeAccessed.String(), info.TimeAccessed.Unix())
// 	}

// 	time.Sleep(1 * time.Second)
// 	cache.Get("b")
// 	for e := cache.list.Front(); e != nil; e = e.Next() {
// 		info := e.Value.(*ElementInfo)
// 		fmt.Printf("k:%s \t v:%s \t t:%s \t %d \t \n", info.Key, info.Value.(string), info.TimeAccessed.String(), info.TimeAccessed.Unix())
// 	}

// }

type ElementInfo struct {
	Key          string
	Value        interface{}
	TimeAccessed time.Time
}

type Cache struct {
	sync.RWMutex
	data map[string]*list.Element
	list *syncList
}

func NewCache(initCapacity int) *Cache {

	return &Cache{
		data: make(map[string]*list.Element, initCapacity),
		list: new(syncList),
	}
}

func (this *Cache) Put(key string, value interface{}) error {
	this.Lock()
	defer this.Unlock()

	if element, ok := this.data[key]; ok {
		element.Value.(*ElementInfo).Value = value

		this.sortedUpdate(element)

	} else {
		element = this.list.PushFront(&ElementInfo{Key: key, Value: value, TimeAccessed: time.Now()})
		this.data[key] = element
	}
	return nil
}

func (this *Cache) Get(key string) interface{} {
	this.RLock()
	defer this.RUnlock()

	if element, ok := this.data[key]; ok {

		this.sortedUpdate(element)
		return element.Value
	}

	return nil
}

func (this *Cache) Delete(key string) error {
	this.Lock()
	defer this.Unlock()

	delete(this.data, key)
	return nil
}

func (this *Cache) Gc(maxlifetime int64) {
	this.Lock()
	defer this.Unlock()

	for {
		element := this.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*ElementInfo).TimeAccessed.Unix() + maxlifetime) < time.Now().Unix() {
			this.list.Remove(element)
			delete(this.data, element.Value.(*ElementInfo).Key)
		} else {
			break
		}
	}
}

func (this *Cache) sortedUpdate(element *list.Element) {
	element.Value.(*ElementInfo).TimeAccessed = time.Now()
	this.list.MoveToFront(element)
}

type syncList struct {
	list.List
	lock sync.Mutex
}

func (this *syncList) MoveToFront(element *list.Element) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.List.MoveToFront(element)
}
