package spmc

import "testing"
import "sync"
import "fmt"

func TestQueue(t *testing.T) {
	var wg sync.WaitGroup
	q := New()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				q.Push(i*100 + j)
			}
		}(i)
	}

	wg.Wait()

	c := 0
	for {
		fmt.Println(q.Pop())
		c++
		if c == 100*100 {
			return
		}
	}
}
