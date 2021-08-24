package lru

import (
	"reflect"
	"testing"
)

type String string

func (o String) Len() int {
	return len(o)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("a", String("111"))

	if v, ok := lru.Get("a"); !ok || v.(String) != String("111") {
		t.Fatal("lru hit key=a value=111 failed")
	}
	if _, ok := lru.Get("b"); ok {
		t.Fatal("lru miss key=b failed")
	}
}

func TestAdd(t *testing.T) {
	lru := New(int64(10), nil)
	lru.Add("a", String("1111"))
	lru.Add("b", String("2222"))
	lru.Add("c", String("3333"))

	if _, ok := lru.Get("a"); ok {
		t.Fatal("key=a should not exists")
	}
	if _, ok := lru.Get("b"); !ok {
		t.Fatal("key=b should exists")
	}
	if _, ok := lru.Get("c"); !ok {
		t.Fatal("key=c should exists")
	}
	if lru.curBytes != int64(len("b") + String("2222").Len() +
		len("c") + String("3333").Len()) {
		t.Fatal("curBytes should be 10")
	}
}

func TestRemoveOldest(t *testing.T) {
	lru := New(int64(10), nil)
	lru.Add("a", String("1111"))
	lru.Add("b", String("2222"))
	lru.RemoveOldest()

	if _, ok := lru.Get("a"); ok {
		t.Fatal("key=a should not exists")
	}
	if _, ok := lru.Get("b"); !ok {
		t.Fatal("key=b should exists")
	}
}

func TestOnEvicted(t *testing.T) {
	var keys []string
	callback := func(key string, val Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("a", String("1111"))
	lru.Add("b", String("2222"))
	lru.Add("c", String("3333"))
	lru.Add("d", String("4444"))

	expect := []string{"a", "b"}
	if !reflect.DeepEqual(expect, keys) {
		t.Fatal("call OnEvicted failed")
	}
}
