package group

import (
	"Geecache/Geecache"
	"Geecache/byteview"
	"Geecache/peerpicker"
	"sync"
)

// 回调接口
type Getter interface {
	Get(key string) ([]byte, error)
}

// 定义一个接口型函数
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// 缓存空间
type Group struct {
	name   string
	getter Getter
	cache  *Geecache.GeeCache
	peers  peerpicker.PeerPicker
}

var (
	Groups      = make(map[string]*Group)
	GroupsGuard = sync.RWMutex{}
)

func NewGroup(name string, cacheByte int, hook Getter) *Group {
	if hook == nil {
		panic("nil hook")
	}

	g := &Group{
		name:   name,
		getter: hook,
		cache:  new(Geecache.GeeCache).SetMaxByte(cacheByte),
	}

	GroupsGuard.Lock()
	defer GroupsGuard.Unlock()

	Groups[name] = g
	return g
}

// 注册一个P2P
func (g *Group) RegisterPeers(peers peerpicker.PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

// 根据组名获取相关指针
func GetGroup(name string) *Group {
	GroupsGuard.RLock()
	defer GroupsGuard.RUnlock()

	return Groups[name]
}

// 从本地加载key的val
func (g *Group) loadFromLocal(key string) (b byteview.Byteview, e error) {
	bs, e := g.getter.Get(key)
	if e != nil {
		return b, e
	}
	b.Init(bs)
	g.populateCache(key, b)
	return b, nil
}

func (g *Group) loadFromPeer(getter peerpicker.PeerGetter, key string) (b byteview.Byteview, e error) {
	if data, err := getter.Get(g.name, key); err == nil {
		b.Init(data)
		return b, nil
	} else {
		return b, err
	}

}

// 加载到cache里面
func (g *Group) populateCache(key string, b byteview.Byteview) {
	g.cache.Add(key, b)
}

func (g *Group) load(key string) (b byteview.Byteview, e error) {
	if g.peers != nil {
		if peer, ok := g.peers.PickPeer(key); ok {
			return g.loadFromPeer(peer, key)
		}
	}
	return g.loadFromLocal(key)
}

// 如果在缓存中则直接取出来，如果不在缓存中，则加载到内存中后返回相关的值
func (g *Group) Get(key string) (b byteview.Byteview, e error) {
	if b, ok := g.cache.Get(key); ok {
		return b, nil
	}

	return g.load(key)
}
