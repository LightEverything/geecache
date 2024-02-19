package p2p_server

import (
	"Geecache/consistenthash"
	"Geecache/peerpicker"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

const DefaultBasePath = "/_geecache"
const DefaultReplicas = 50

type HTTPPool struct {
	self       string
	basePath   string
	peers      *consistenthash.Map
	peersGuard sync.Mutex
	httpGetter map[string]*httpGetter
}

type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	// 获取网址
	u := fmt.Sprintf(
		"%v//%v//%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key))

	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("server returned " + res.Status)
	}

	// 获取body数据
	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: DefaultBasePath,
	}
}

// 参数设置Peer
func (h *HTTPPool) Set(peers ...string) {
	h.peersGuard.Lock()
	defer h.peersGuard.Unlock()

	// 一致性hash
	h.peers = consistenthash.New(DefaultReplicas, nil)
	h.peers.Add(peers...)
	h.httpGetter = make(map[string]*httpGetter, len(peers))

	for _, peer := range peers {
		h.httpGetter[peer] = &httpGetter{baseURL: peer + h.basePath}
	}
}

func (h *HTTPPool) PickPeer(key string) (pg peerpicker.PeerGetter, ok bool) {
	h.peersGuard.Lock()
	defer h.peersGuard.Unlock()

	if peer := h.peers.Get(key); peer != "" && peer != h.self {

		return h.httpGetter[peer], true
	}
	return nil, false
}
