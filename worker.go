package gpool

import "runtime"

func newPoolWorker() *poolWorker {
	pw := &poolWorker{
		wc: make(chan func(), 8),
	}

	go pw.loop()
	runtime.SetFinalizer(pw, func(w *poolWorker) { w.wc <- nil })
	return pw
}

// work in go env
type poolWorker struct {
	wc chan func()
}

func (w *poolWorker) work(fn func()) {
	w.wc <- fn
}

func (w *poolWorker) loop() {
	for fn := range w.wc {
		if fn == nil {
			return
		}

		fn()
	}
}
