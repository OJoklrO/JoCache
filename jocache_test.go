package JoCache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(s string) ([]byte, error) {
		return []byte(s), nil
	})

	expect := []byte("key")

	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Fatal("error")
	}
}

var db = map[string]string {
	"Joklr": "123",
	"abc": "#sdfs",
	"asd": "xzc",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	j := NewGroup("test", 2 << 10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key]++
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		if view, err := j.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value")
		}
		if _, err := j.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}

	if _, err := j.Get("unknown"); err == nil {
		t.Fatal("asdasda")
	}
}