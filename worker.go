package gpool

import "runtime"

func newPoolWorker() *poolWorker {
	pw := &poolWorker{
		wc: make(chan *worker),
	}

	go pw.loop()
	runtime.SetFinalizer(pw, func(w *poolWorker) { w.wc <- nil })
	return pw
}

// work in go env
type poolWorker struct {
	wc chan *worker
}

type worker struct {
	fn      func() error
	retChan chan error
}

func (w *poolWorker) work(fn func() error) error {
	ret := make(chan error)
	w.wc <- &worker{
		fn:      fn,
		retChan: ret,
	}
	err := <-ret
	return err
}

func (w *poolWorker) loop() {
	for ch := range w.wc {
		if ch == nil {
			return
		}

		ch.retChan <- ch.fn()
	}
}
