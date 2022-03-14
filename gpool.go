package gpool

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	gp = NewGoPool(10000)
)

// Go run fn
func Go(fn func()) {
	gp.Go(fn)
}

var GoPoolMaxGoroutineError = errors.New("reach max goroutine count.")

type GoPoolOptions struct {
	reject func(fn func()) error
}

func (o *GoPoolOptions) apply() {
	if o.reject == nil {
		o.reject = func(func()) error {
			return GoPoolMaxGoroutineError
		}
	}
}

type GoPoolOption func(*GoPoolOptions)

func WithReject(fn func(func()) error) GoPoolOption {
	return func(options *GoPoolOptions) {
		options.reject = fn
	}
}

func NewGoPool(maxCount int32, opts ...GoPoolOption) *GoPool {
	var opt GoPoolOptions
	for _, o := range opts {
		o(&opt)
	}

	opt.apply()
	return &GoPool{
		opts: opt,
		poolWorker: sync.Pool{
			New: func() interface{} {
				return newPoolWorker()
			},
		},
		maxCount: maxCount,
	}
}

// GoPool go pool
type GoPool struct {
	poolWorker sync.Pool

	maxCount int32
	curCount int32
	opts     GoPoolOptions
}

// Go pool will try get a goroutine from pool to exec fn
// if goroutine count reach max count, GoPool will use reject strategy handle.
func (p *GoPool) Go(fn func()) error {
	if atomic.AddInt32(&p.curCount, 1) > p.maxCount {
		atomic.AddInt32(&p.curCount, -1)
		return p.opts.reject(fn)
	}
	pw := p.poolWorker.Get().(*poolWorker)
	pw.work(func() {
		fn()
		atomic.AddInt32(&p.curCount, -1)
		p.poolWorker.Put(pw)
	})
	return nil
}

// RunningTasks running tasks count
func (p *GoPool) RunningTasks() int32 {
	return atomic.LoadInt32(&p.curCount)
}
