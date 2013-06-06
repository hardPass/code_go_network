package concurrent

import (
  "fmt"
  // "runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// This program shows how a `copy-on-write` mode works.

type CopyOnWriteSlice struct {
	ptr_snapshot unsafe.Pointer
	valid        chan bool
	closeLock    sync.Mutex
	rmFilter     func(interface{}) bool
}

func New(cap uint32) *CopyOnWriteSlice {
	tmp := make([]*interface{}, 0, cap)
	ptr_snapshot = unsafe.Pointer(&tmp)

	return &CopyOnWriteSlice{
		ptr_snapshot: ptr_snapshot,
		valid:        make(chan bool),
	}
}

func (c *CopyOnWriteSlice) RemoveFilterFunc(filter func(interface{}) bool) *CopyOnWriteSlice {
	c.rmFilter = filter
}

func (c *CopyOnWriteSlice) Start(clean_interval time.Duration) *CopyOnWriteSlice {
	c.rmFilter = filter
}

func (c *CopyOnWriteSlice) Close() {
	c.closeLock.Lock()
	defer c.closeLock.Unlock()

	if c.Closed() {
		return
	}

	close(c.valid)
}

// Did you see the trick here?
func (c *CopyOnWriteSlice) Closed() bool {
	select {
	case _, ok := <-c.valid:
		if !ok {
			return true
		}
	default:
	}

	return false
}
