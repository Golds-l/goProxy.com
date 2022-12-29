package buffer

import (
	"fmt"
	"io"
	"sync"
	"testing"
	"time"
)

type testW int

func (t testW) Write(p []byte) (n int, err error) {
	fmt.Print(p)
	return 0, nil
}

func (b *Buffer) Print() {
	fmt.Println()
	var n = b.R.n
	for {
		if n == b.R.n && n == b.W.n {
			fmt.Println(n.Content, "Read&Write pointer  read offset:", b.R.offset, "write offset:", b.W.offset)
		} else if n == b.R.n {
			fmt.Println(n.Content, "Read pointer  offset:", b.R.offset)
		} else if n == b.W.n {
			fmt.Println(n.Content, "Write pointer  offset:", b.W.offset)
		} else {
			fmt.Println(n.Content)
		}
		if n.Next == b.R.n {
			break
		}
		n = n.Next
	}
}

func (b *Buffer) Fill(s, e int) {
	bytes := make([]byte, e-s)
	for i := s; s < e; s++ {
		bytes[s-i] = uint8(s)
	}
	b.Write(bytes)
}

func (b *Buffer) GoFill(s, e int, wg *sync.WaitGroup) {
	for ; s < e; s++ {
		b.Fill(s, s+1)
		//time.Sleep(1 * time.Second)
		wg.Done()
	}
}

func (b *Buffer) GoRead(n int, w io.Writer, wg *sync.WaitGroup) {
	for n > 0 {
		b.Read(2, w)
		n -= 2
		time.Sleep(1 * time.Second)
		wg.Done()
	}
}

func TestBuffer_Write(t *testing.T) {
	b := &Buffer{}
	wg := &sync.WaitGroup{}
	var w testW
	b.InitBuffer(10, 10)
	b.Fill(0, 10)
	wg.Add(148)
	go b.GoFill(10, 120, wg)
	go b.GoRead(56, w, wg)
	wg.Wait()
	b.Print()
}
