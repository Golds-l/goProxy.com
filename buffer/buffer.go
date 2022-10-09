package buffer

import "sync"

type Buffer struct {
	W       *Node
	R       *Node
	NodeNum int
	NodeLen int
	use     int
	mu      sync.Mutex
}

type Node struct {
	Next    *Node
	Pre     *Node
	Content []byte
	Len     int
}

func makeListNode(n, l int, head *Node) *Node {
	var node = &Node{Pre: head, Next: nil, Content: make([]byte, l), Len: l}
	for i := 0; i < n; i++ {
		var t = Node{Pre: node, Next: nil, Content: make([]byte, l), Len: l}
		node.Next = &t
		node = &t
	}
	return node
}

func (b *Buffer) InitBuffer(n int, l int) {
	b.W = &Node{Pre: nil, Next: nil, Content: nil, Len: 0}
	b.R = nil
	h := makeListNode(n, l, b.W)
	h.Next = b.W
	b.R = b.W
}
