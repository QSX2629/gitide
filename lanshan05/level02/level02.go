package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Task struct {
	Adder func(an int) int64
}

func main() {
	ch := make(chan Task, 10)
	wg := sync.WaitGroup{}
	var totalsum int64
	for an := range 10 {
		wg.Add(1)
		go func(an int) {
			defer wg.Done()
			for t := range ch {
				re := t.Adder(an)
				atomic.AddInt64(&totalsum, re)
			}
		}(an)
	}
	for i := 0; i < 20; i++ {
		j := i
		t1 := Task{
			Adder: func(an int) int64 {
				return int64(j + 1)
			},
		}
		ch <- t1
	}
	close(ch)
	wg.Wait()
	fmt.Println(totalsum)
	finalsum := atomic.LoadInt64(&totalsum)
	fmt.Println(finalsum)
}
