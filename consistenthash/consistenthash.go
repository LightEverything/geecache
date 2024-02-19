package consistenthash

import (
	"Geecache/utility"
	"hash/crc32"
	"slices"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash Hash
	// 虚拟节点个数
	replicas int
	keys     []int // sorted
	hashMap  map[int]string
}

func New(replicas int, f Hash) *Map {
	m := &Map{
		hash:     f,
		replicas: replicas,
		keys:     make([]int, 0),
		hashMap:  make(map[int]string),
	}

	if f == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

// 节点名称
func (m *Map) Add(keys ...string) {
	for _, k := range keys {
		for i := 0; i < m.replicas; i++ {
			h := int(m.hash([]byte(strconv.Itoa(i) + k)))
			idx := utility.UpperBoundInt(m.keys, h)
			m.keys = slices.Insert(m.keys, idx, h)
			m.hashMap[h] = k
		}
	}
}

// 选择节点
func (m *Map) Get(key string) (node string) {
	if len(m.keys) == 0 {
		return ""
	}

	h := int(m.hash([]byte(key)))

	idx := utility.LowerBoundInt(m.keys, h)
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
