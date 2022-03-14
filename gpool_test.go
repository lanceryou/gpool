package gpool

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGoPool_Go(t *testing.T) {
	ts := []struct {
		maxCount int32
		cnt      int64
		expect   int64
	}{
		{
			maxCount: 0,
		},
		{
			maxCount: 1,
			expect:   1,
		},
	}

	for _, s := range ts {
		p := NewGoPool(s.maxCount)
		for i := 0; i < int(s.maxCount)+1; i++ {
			err := p.Go(func() {
				atomic.AddInt64(&s.cnt, 1)
				time.Sleep(time.Second)
			})

			if i < int(s.maxCount) && err != nil {
				t.Errorf("expect err nil, but err %v", err)
			} else if i >= int(s.maxCount) && err != GoPoolMaxGoroutineError {
				t.Errorf("i %v expect %v, but error %v", i, GoPoolMaxGoroutineError, err)
			}
		}
	}
}

func TestNewGoPool(t *testing.T) {
	ts := []struct {
		maxCount int32
		cnt      int64
		expect   int64
	}{
		{
			maxCount: 0,
		},
		{
			maxCount: 1,
			expect:   1,
		},
	}

	for _, s := range ts {
		var wg sync.WaitGroup
		p := NewGoPool(s.maxCount)
		for i := 0; i < int(s.maxCount)+1; i++ {
			err := p.Go(func() {
				atomic.AddInt64(&s.cnt, 1)
				time.Sleep(time.Second)
				wg.Done()
			})
			if err == nil {
				wg.Add(1)
			}

			if i < int(s.maxCount) && err != nil {
				t.Errorf("expect err nil, but err %v", err)
			} else if i >= int(s.maxCount) && err != GoPoolMaxGoroutineError {
				t.Errorf("i %v expect %v, but error %v", i, GoPoolMaxGoroutineError, err)
			}
		}

		wg.Wait()
		if s.cnt != s.expect {
			t.Errorf("expect %v, but %v", s.expect, s.cnt)
		}
	}
}

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
)

func TestNoPool(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10000000; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Millisecond * 10)
			wg.Done()
		}()
	}

	wg.Wait()
	var curMem uint64
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}

func TestGoPool(t *testing.T) {
	var wg sync.WaitGroup

	p := NewGoPool(10000000)

	for i := 0; i < 10000000; i++ {
		wg.Add(1)
		p.Go(func() {
			time.Sleep(time.Millisecond * 10)
		})
	}
	var curMem uint64
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem = mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)
}

func BenchmarkGoPool_Go(b *testing.B) {
	ts := []struct {
		maxCount int32
	}{
		{
			maxCount: 1000000,
		},
	}

	for _, s := range ts {
		p := NewGoPool(s.maxCount)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			p.Go(func() {
				time.Sleep(time.Millisecond * 10)
			})
		}
		b.StopTimer()
	}
}

func BenchmarkGoPool_Go2(b *testing.B) {
	p := NewGoPool(10000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Go(func() {
				time.Sleep(time.Millisecond * 10)
			})
		}
	})
}

func BenchmarkGo(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Millisecond * 10)
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkGo2(b *testing.B) {
	var wg sync.WaitGroup
	p := NewGoPool(100000)
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		p.Go(func() {
			time.Sleep(time.Millisecond * 10)
			wg.Done()
		})
	}
	wg.Wait()
}
