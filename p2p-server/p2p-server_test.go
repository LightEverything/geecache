package p2p_server

import (
	"Geecache/group"
	"fmt"
	"log"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGetRouter(t *testing.T) {

	group.NewGroup("scores", 2<<10, group.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	r := GetRouter()
	log.Println("geecache is running at", addr)
	log.Fatal(r.Run(addr))
}
