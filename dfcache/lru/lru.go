package lru

import "container/list"

type LRU struct {
	maxBytes  int64
	curBytes  int64
	linkList  *list.List
	linkMap   map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type Value interface {
	Len() int
}

// The inner struct of the LRU
type entry struct {
	key string
	val Value
}

// Add item to the LRU
func (c *LRU) Add(key string, val Value) {
	if ele, ok := c.linkMap[key]; ok {
		c.linkList.MoveToFront(ele)
		en := ele.Value.(*entry)
		c.curBytes += int64(val.Len() - en.val.Len())
		en.val = val
	} else {
		ele := c.linkList.PushFront(&entry{key, val})
		c.linkMap[key] = ele
		c.curBytes += int64(len(key) + val.Len())
	}

	// Check maxBytes and delete the oldest item
	for c.maxBytes != 0 && c.curBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

// Remove the oldest item from the LRU
func (c *LRU) RemoveOldest() {
	ele := c.linkList.Back()
	if ele == nil {
		return
	}

	c.linkList.Remove(ele)
	en := ele.Value.(*entry)
	delete(c.linkMap, en.key)
	c.curBytes -= int64(len(en.key) + en.val.Len())

	if c.OnEvicted == nil {
		return
	}
	c.OnEvicted(en.key, en.val)
}

// Get item from the LRU
func (c *LRU) Get(key string) (val Value, ok bool) {
	if ele, ok := c.linkMap[key]; ok {
		c.linkList.MoveToFront(ele)
		en := ele.Value.(*entry)
		return en.val, ok
	}
	return
}

// Get the number of the linklist entries
func (c *LRU) Len() int {
	return c.linkList.Len()
}

// New the instance of the LRU
func New(maxBytes int64, onEvicted func(string, Value)) *LRU {
	return &LRU{
		maxBytes:  maxBytes,
		linkList:  list.New(),
		linkMap:   make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}
