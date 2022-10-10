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
	WriteOff int
	ReadOff  int
	Next     *Node
	Pre      *Node
	Content  []byte
	Len      int
}

func (head *Node) makeListNode(n, l int) *Node {
	var node = &Node{Pre: head, Next: nil, Content: make([]byte, l), Len: l}
	head.Next = node
	for i := 2; i < n; i++ {
		var t = Node{Pre: node, Next: nil, Content: make([]byte, l), Len: l}
		node.Next = &t
		node = &t
	}
	return node
}

func (b *Buffer) InitBuffer(n int, l int) {
	b.W = &Node{Len: 0, Content: make([]byte, l)}
	tail := (b.W).makeListNode(n, l)
	tail.Next = b.W
	b.W.Pre = tail
	b.R = b.W
	b.use = 0
	b.NodeNum = n
	b.NodeLen = l
}

func (b *Buffer) Write(c []byte) {
	l := len(c)
	p := b.W
	for l > 0 {
		if l < b.NodeLen-p.WriteOff {
			copy(p.Content[p.WriteOff:p.WriteOff+l], c[len(c)-l:])
			p.WriteOff += l
			l = 0
		} else {
			copy(p.Content, c[len(c)-l:len(c)-l+b.NodeLen])
			p.WriteOff = b.NodeLen
			p = p.Next
			l -= b.NodeLen
		}
	}
	b.W = p
}

func (b *Buffer) Read(n int, w io.Writer) {
	p := b.R
	for n > 0 {
		if n > b.NodeLen-p.ReadOff {
			w.Write(p.Content[p.ReadOff:])
			n -= b.NodeLen - p.ReadOff
			p.ReadOff = b.NodeLen
			p = p.Next
		} else {
			w.Write(p.Content[p.ReadOff : p.ReadOff+n])
			p.ReadOff += n
			n = 0
		}
	}
	b.R = p
}

func (b *Buffer) Print() {
	n := b.W
	for {
		fmt.Println(n)
		if n == b.W.Pre {
			break
		}
		n = n.Next
	}
}
