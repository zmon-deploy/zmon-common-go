package misc

import (
	"sync"
	"sync/atomic"
)

type ExpirableOnce struct {
	m    sync.Mutex
	done uint32
}

func (o *ExpirableOnce) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}

	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

func (o *ExpirableOnce) Expire() {
	if atomic.LoadUint32(&o.done) == 0 {
		return
	}

	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 1 {
		atomic.StoreUint32(&o.done, 0)
	}
}
