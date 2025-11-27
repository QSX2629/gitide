package taskpool

import "sync"

type Task interface {
	Execute()
}
type Pool struct {
	ch chan Task
	wg sync.WaitGroup
}

func New(gn int, tbs int) *Pool {
	p := &Pool{
		ch: make(chan Task, tbs),
	}
	p.wg.Add(gn)
	for i := 0; i < gn; i++ {
		go func() {
			defer p.wg.Done()
			for t := range p.ch {
				t.Execute()
			}
		}()
	}
	return p
}
func (p *Pool) Submit(task Task) {
	p.ch <- task
}
func (p *Pool) Close() {
	close(p.ch)
	p.wg.Wait()
}
