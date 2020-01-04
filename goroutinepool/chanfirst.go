package goroutinepool

import "sync/atomic"

type chanFirstPool struct {
	cap           uint64
	runningWorker uint64
	task          chan func()
	close         chan struct{}
}

func NewChanFirstPool(cap uint64) *chanFirstPool {
	return &chanFirstPool{
		cap:  cap,
		task: make(chan func(), cap),
	}
}

func (p *chanFirstPool) Schedule(task func()) {
	if p.getRunningWorker() < p.cap {
		go p.run()
	}
	p.task <- task
}

func (p *chanFirstPool) run() {
	defer p.decrease()

	p.increase()
	for {
		select {
		case task, ok := <-p.task:
			if !ok {
				return
			}
			task()
		case <-p.close:

		}
	}
}

func (p *chanFirstPool) increase() {
	atomic.AddUint64(&p.runningWorker, 1)
}

// a, b > 0
// a - b = a- b + mod = a + (mod - b)
// (mod - b) 的 bit 表示
// mod = 1{n} + 1
// mod - b = 1{n} + 1 - b = 1{n} - (b - 1) = ^(b - 1)
func (p *chanFirstPool) decrease() {
	atomic.AddUint64(&p.runningWorker, ^uint64(0))
}

func (p *chanFirstPool) getRunningWorker() uint64 {
	return atomic.LoadUint64(&p.runningWorker)
}
