package gpool

import (
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
