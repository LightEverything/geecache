package lru

import "container/list"

// 删除调用的Hook函数
type EvictedHandler func(string, Value)

// 缓存结构体
type Cache struct {
	ll       *list.List
	cacheMap map[string]*list.Element
	maxByte  int
	nByte    int

	onEvicted EvictedHandler
}

type entry struct {
	key   string
	value Value
}

// 存储值的struct
type Value interface {
	Len() int
}

func (c *Cache) New(maxByte int, onEvicted EvictedHandler) *Cache {
	return &Cache{
		ll:        new(list.List),
		cacheMap:  make(map[string]*list.Element),
		maxByte:   maxByte,
		nByte:     0,
		onEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (v Value, ok bool) {
	if node, ok := c.cacheMap[key]; ok {
		// lru缓存策略，最新访问的放在最前面
		c.ll.MoveToFront(node)
		// 返回entry结构体
		return node.Value.(*entry).value, true
	} else {
		return v, false
	}
}

func (c *Cache) Add(key string, v Value) {
	if ele, ok := c.cacheMap[key]; ok {
		// 如果存在，则修改并且更新策略
		c.ll.MoveToFront(ele)

		// 修改原有值
		tmp := ele.Value.(*entry).value
		ele.Value.(*entry).value = v

		c.nByte += v.Len() - tmp.Len()
	} else {
		node := c.ll.PushFront(&entry{key, v})
		c.cacheMap[key] = node
		c.nByte += v.Len() + len(key)
	}

	// 如果超出最大限制,则更新缓存
	if c.nByte > c.maxByte && c.maxByte != 0 {
		c.RemoveOldest()
	}
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		// 获取ele 元素
		tmp := ele.Value.(*entry)
		c.ll.Remove(ele)
		delete(c.cacheMap, tmp.key)
		c.nByte -= len(tmp.key) + tmp.value.Len()

		// 调用HOOK
		if c.onEvicted != nil {
			c.onEvicted(tmp.key, tmp.value)
		}
	}
}

// 返回存储数据
func (c *Cache) Len() int {
	return c.nByte
}
