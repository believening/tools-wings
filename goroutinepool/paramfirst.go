package goroutinepool

type paramFirstPool struct {
	work chan func()
	sem  chan struct{}
}

func NewParamFirst(size int) Pool {
	return &paramFirstPool{
		work: make(chan func()),
		sem:  make(chan struct{}, size),
	}
}

func (p *paramFirstPool) Schedule(task func()) {
	select {
	case p.work <- task:
	case p.sem <- struct{}{}:
		go p.run(task)
	}
}

func (p *paramFirstPool) run(task func()) {
	defer func() { <-p.sem }()
	for {
		task()
		task = <-p.work // 复用 goroutine
	}
}
