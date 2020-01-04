package goroutinepool

type Pool interface {
	Schedule(func())
}
