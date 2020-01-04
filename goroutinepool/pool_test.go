package goroutinepool

import (
	"sync"
	"sync/atomic"
	"testing"
)

var (
	wg         = sync.WaitGroup{}
	benchTimes = 1000000
	sum        int64
)

func addTask() {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		atomic.AddInt64(&sum, 1)
	}
}

func BenchmarkGoroutines(b *testing.B) {
	for i := 0; i < benchTimes; i++ {
		wg.Add(1)
		go addTask()
	}
	wg.Wait()
}

func BenchmarkParamFirstPool(b *testing.B) {
	p := NewParamFirst(20)
	b.ResetTimer()
	for i := 0; i < benchTimes; i++ {
		wg.Add(1)
		p.Schedule(addTask)
	}
	wg.Wait()
}

func BenchmarkChanFirstPool(b *testing.B) {
	p := NewChanFirstPool(20)
	b.ResetTimer()
	for i := 0; i < benchTimes; i++ {
		wg.Add(1)
		p.Schedule(addTask)
	}
	wg.Wait()
}
