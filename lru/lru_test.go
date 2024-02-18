package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

var testUnits = map[string]String{
	"key1": "v1",
	"key2": "v2",
	"key3": "v3",
}

func TestGet(t *testing.T) {
	testLru := (&Cache{}).New(20, nil)

	for k, v := range testUnits {
		testLru.Add(k, v)
		if r, ok := testLru.Get(k); !ok || r.(String) != v {
			t.Fatal()
		}
	}
}

func TestRemoveoldest(t *testing.T) {
	testLru := (&Cache{}).New(10, nil)

	for k, v := range testUnits {
		testLru.Add(k, v)
	}

	if _, ok := testLru.Get("key1"); ok {
		t.Fatalf("Removeoldest key1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	testLru := (&Cache{}).New(10, callback)
	testLru.Add("key1", String("123456"))
	testLru.Add("k2", String("k2"))
	testLru.Add("k3", String("k3"))
	testLru.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}

func TestAdd(t *testing.T) {
	testLru := (&Cache{}).New(6, nil)
	testLru.Add("key", String("1"))
	//log.Println(testLru.nByte)
	testLru.Add("key", String("111"))
	//log.Println(testLru.nByte)
	v, ok := testLru.Get("key")

	if testLru.Len() != len("key")+len("111") || !ok || v != String("111") {
		t.Fatal(ok)
	}
}
