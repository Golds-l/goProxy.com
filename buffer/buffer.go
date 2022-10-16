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
	mu        sync.Mutex
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
}

func (b *Buffer) Write(c []byte) {
	l := len(c)
	p := b.W.n
	for l > 0 {
		writePointer := b.W
		// TODO: 读写锁需要分离
		// 锁定当前读指针
		b.mu.Lock()
		readPointer := b.R
		b.mu.Unlock()
		// 临界节点写，读写指针处于同一节点
		if writePointer.n == readPointer.n && readPointer.offset != writePointer.offset {
			// 剩余容量不足，等待
			if writePointer.offset+l >= readPointer.offset {
				b.writeCond.Wait()
				continue
			} else {
				copy(p.Content[writePointer.offset:writePointer.offset+l], c[len(c)-l:])
				writePointer.offset += l
				l = 0
			}
		} else {
			// 当前节点足够写
			if l < b.NodeLen-writePointer.offset {
				copy(p.Content[writePointer.offset:writePointer.offset+l], c[len(c)-l:])
				writePointer.offset += l
				l = 0
				// 当前节点不够写，全部写满
			} else {
				copy(p.Content[writePointer.offset:], c[len(c)-l:len(c)-l+b.NodeLen-writePointer.offset])
				l -= b.NodeLen - writePointer.offset
				p = p.Next
				writePointer.n = p
				writePointer.offset = 0
			}
		}
	}
	b.W.n = p
}

func (b *Buffer) Read(n int, w io.Writer) {
	p := b.R.n
	for n > 0 {
		// 临界节点读 锁定当前写指针
		b.mu.Lock()
		writePointer := b.W
		b.mu.Unlock()
		readPointer := b.R
		if b.R.n == b.W.n {
			// 当前内容足够读长度
			if writePointer.offset-readPointer.offset >= n {
				_, _ = w.Write(p.Content[readPointer.offset : readPointer.offset+n])
				readPointer.offset += n
				n = 0
			} else if writePointer.offset == readPointer.offset {
				// 当前读指针与写指针处于同一位置，等待
				b.readCond.Wait()
				continue
			} else {
				// 当前不足够读，有多少读多少
				_, _ = w.Write(p.Content[readPointer.offset:writePointer.offset])
				readPointer.offset = writePointer.offset
				n -= writePointer.offset - readPointer.offset
			}
		} else {
			// 未到临界点，足够读
			if n > b.NodeLen-readPointer.offset {
				// 当前节点内容全部读入
				_, _ = w.Write(p.Content[readPointer.offset:])
				n -= b.NodeLen - readPointer.offset
				p = p.Next
				readPointer.n = p
				readPointer.offset = 0
			} else {
				// 当前节点部分读入
				_, _ = w.Write(p.Content[readPointer.offset : readPointer.offset+n])
				readPointer.offset += n
				n = 0
			}
		}
	}
	b.R.n = p
}
