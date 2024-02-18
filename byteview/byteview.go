package byteview

// 只读的byte,防止外部程序修改cache(注意：不安全的）
type Byteview struct {
	b      []byte
	isInit bool
}

// 返回byteview的长度,实现value的接口
func (b Byteview) Len() int {
	return len(b.b)
}

func (b *Byteview) Init(bs []byte) {
	if b.isInit == false {
		b.isInit = true
		b.b = make([]byte, len(bs))
		copy(b.b, bs)
	} else {
		panic("init byteview more")
	}
}

// 复制切片，因为byteview是只读
func (b Byteview) ByteSlice() (d []byte) {
	d = make([]byte, len(b.b))
	copy(d, b.b)
	return d
}

func (b Byteview) String() string {
	return string(b.b)
}
