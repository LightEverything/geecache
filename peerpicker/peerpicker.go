package peerpicker

type PeerPicker interface {
	PickPeer(key string) (pg PeerGetter, ok bool)
}

type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
