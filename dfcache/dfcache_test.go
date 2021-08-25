package dfcache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	key := "aaa"
	expect := []byte(key)
	if v, _ := f.Get(key); !reflect.DeepEqual(v, expect) {
		t.Fatal("getter failed")
	}
}

func TestGetGroup(t *testing.T) {
	gName := "dfGroup"
	NewGroup(gName, 2<<10, GetterFunc(
		func(key string) (bytes []byte, err error) { return }))

	if g := GetGroup(gName); g == nil || g.name != gName {
		t.Fatal("GetGroup failed")
	}

	if g := GetGroup(gName + "111"); g != nil {
		t.Fatal("GetGroup failed, un exists")
	}
}

func TestGet(t *testing.T) {
	var db = map[string]string{
		"a": "111",
		"b": "222",
		"c": "333",
	}
	loadCounts := make(map[string]int, len(db))

	dfcache := NewGroup("dfcache", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[db] search key", key)
			if v, ok := db[key]; ok {
				loadCounts[key]++
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exists", key)
		}))

	for k, v := range db {
		if view, err := dfcache.Get(k); err != nil || view.String() != v {
			t.Fatalf("failed to get key of %s", k)
		}
		if _, err := dfcache.Get(k); err != nil || loadCounts[k] != 1 {
			t.Fatalf("cache %s miss", k)
		}
	}

	if _, err := dfcache.Get("unknow"); err == nil {
		t.Fatalf("key %s should not exists", "unknow")
	}
}
