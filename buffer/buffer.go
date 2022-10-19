package buffer

import (
	"io"
	"sync"
)

type Node struct {
	Pre     *Node
	Next    *Node
	Content []byte
	Len     int
}

type pointer struct {
	n      *Node
	offset int
}

type Buffer struct {
	W         pointer
	R         pointer
	NodeNum   int
	NodeLen   int
	readMu    *sync.Mutex
	writeMu   *sync.Mutex
	writeCond *sync.Cond
	readCond  *sync.Cond
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
	b.W =
		pointer{
			&Node{Len: 0, Content: make([]byte, l)},
			0,
		}
	tail := (b.W.n).makeListNode(n, l)
	tail.Next = b.W.n
	b.W.n.Pre = tail
	b.R = b.W
	b.NodeNum = n
	b.NodeLen = l
	b.readCond = sync.NewCond(&sync.Mutex{})
	b.writeCond = sync.NewCond(&sync.Mutex{})
	b.readMu = &sync.Mutex{}
	b.writeMu = &sync.Mutex{}
}

func (b *Buffer) Write(c []byte) {
	l := len(c)
	p := b.W.n
	//b.W := b.W
	for l > 0 {
		// 锁定当前读指针
		b.readMu.Lock()
		readPointer := b.R
		b.readMu.Unlock()
		// 临界节点写，读写指针处于同一节点
		if b.W.n == readPointer.n && readPointer.offset > b.W.offset {
			// 剩余容量不足，等待
			if b.W.offset+l >= readPointer.offset {
				b.writeCond.L.Lock()
				b.writeCond.Wait()
				b.writeCond.L.Unlock()
				continue
			} else {
				copy(p.Content[b.W.offset:b.W.offset+l], c[len(c)-l:])
				b.W.offset += l
				l = 0
				b.readCond.Broadcast()
			}
		} else {
			// 当前节点足够写
			if l < b.NodeLen-b.W.offset {
				copy(p.Content[b.W.offset:b.W.offset+l], c[len(c)-l:])
				b.W.offset += l
				l = 0
				b.readCond.Broadcast()
			} else {
				// 当前节点不够写，写满该节点
				copy(p.Content[b.W.offset:], c[len(c)-l:len(c)-l+b.NodeLen-b.W.offset])
				l -= b.NodeLen - b.W.offset
				p = p.Next
				b.W.n = p
				b.W.offset = 0
				b.readCond.Broadcast()
			}
		}
	}
	b.W = b.W
}

func (b *Buffer) Read(n int, w io.Writer) {
	p := b.R.n
	//b.R := b.R
	for n > 0 {
		// 锁定当前写指针
		//b.writeMu.Lock()
		//b.W := b.W
		//b.writeMu.Unlock()
		// 临界节点读
		if b.R.n == b.W.n {
			// 当前内容足够读长度
			if b.W.offset-b.R.offset >= n {
				_, _ = w.Write(p.Content[b.R.offset : b.R.offset+n])
				b.R.offset += n
				n = 0
				b.writeCond.Broadcast()
			} else if b.W.offset == b.R.offset {
				// 当前读指针与写指针处于同一位置，等待
				b.readCond.L.Lock()
				b.readCond.Wait()
				b.readCond.L.Unlock()
				continue
			} else {
				// 当前不足够读，有多少读多少
				// TODO: b.W 指针未锁定
				_, _ = w.Write(p.Content[b.R.offset:b.W.offset])
				b.R.offset = b.W.offset
				n -= b.W.offset - b.R.offset
				b.writeCond.Broadcast()
			}
		} else {
			// 未到临界点，足够读
			if n >= b.NodeLen-b.R.offset {
				// 当前节点内容全部读入
				_, _ = w.Write(p.Content[b.R.offset:])
				n -= b.NodeLen - b.R.offset
				p = p.Next
				b.R.n = p
				b.R.offset = 0
				b.writeCond.Broadcast()
			} else {
				// 当前节点部分读入
				_, _ = w.Write(p.Content[b.R.offset : b.R.offset+n])
				b.R.offset += n
				n = 0
				b.writeCond.Broadcast()
			}
		}
	}
	b.R = b.R
}
