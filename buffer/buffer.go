package buffer

type Buffer struct {
	W           *Node
	R           *Node
	NodeNum     int
	NodeDataLen int
}

type Node struct {
	Next       *Node
	Pre        *Node
	Content    []byte
	ContentLen int
}

func (b *Buffer) InitBuffer(n int, l int) {
	b.W = &Node{Pre: nil, Next: nil, Content: nil, ContentLen: 0}
	b.R = nil
	h := b.W
	for i := 0; i < n; i++ {
		var t = Node{Pre: h, Next: nil, Content: nil, ContentLen: i + 1}
		var c = make([]byte, l)
		t.Content = c
		h.Next = &t
		h = &t
		b.NodeNum += 1
	}
	h.Next = b.W
	b.R = b.W
}
