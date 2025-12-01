package taskpool

import (
	"fmt"
	"sync"
)

type Task interface {
	Execute()
}
type Pool struct {
	ch     chan Task
	wg     sync.WaitGroup
	closed bool
	mu     sync.Mutex
}

func New(gn int, tbs int) *Pool {
	if gn <= 0 {
		gn = 5
	}
	if tbs <= 0 {
		tbs = 0
	}
	p := &Pool{
		ch: make(chan Task, tbs),
	}

	p.wg.Add(gn)
	for i := 0; i < gn; i++ {
		go p.Worker()
	}
	return p
}
func (p *Pool) Worker() {
	defer p.wg.Done()
	for task := range p.ch {
		task.Execute()
	}
}
func (p *Pool) Submit(task Task) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		fmt.Println("pool is closed")
		return
	}
	p.ch <- task
}
func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return
	}
	close(p.ch)
	p.wg.Wait()
	p.closed = true
}
