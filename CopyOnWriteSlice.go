package concurrent

import (
	"log"
	// "sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// This program shows how a `copy-on-write` mode works.
// concurrent read: lock free
// concurrent wirte: in `copy-on-write` mode, and all writings operted in one goroutine
// use this in the case like this: reading opertions are much more than writing.

type CopyOnWriteSlice struct {
	ptr_snapshot unsafe.Pointer   //  *[]interface{} ,
	addChan      chan interface{} // elements to append
	cleanMark    chan bool        // a flag for marking to clean by filter func
	closed       chan bool
	filter       func(interface{}) bool // for clean, delete it if return false; otherwise, it will be reserved if it return true
}

func New(cap uint32, filter func(interface{}) bool, clean_interval time.Duration) *CopyOnWriteSlice {
	base := make([]interface{}, 0, cap)

	return Wrap(base, filter, clean_interval)
}

func Wrap(base []interface{}, filterFunc func(interface{}) bool, clean_interval time.Duration) *CopyOnWriteSlice {
	c := &CopyOnWriteSlice{
		ptr_snapshot: unsafe.Pointer(&base),
		addChan:      make(chan interface{}, 10),
		cleanMark:    make(chan bool, 1),
		closed:       make(chan bool, 1),
		filter:       filterFunc,
	}

	go c.loopWrite(clean_interval)

	return c
}

// don't execute this twice
func (c *CopyOnWriteSlice) Close() {
	close(c.closed)
}

// Did you see the trick here?
func (c *CopyOnWriteSlice) Closed() bool {
	select {
	case _, ok := <-c.closed:
		if !ok {
			return true
		}
	default:
	}

	return false
}

// only for read
func (c *CopyOnWriteSlice) Len() int {
	return len(c.Snapshot())
}

// only for read
func (c *CopyOnWriteSlice) Snapshot() []interface{} {
	p := atomic.LoadPointer(&c.ptr_snapshot)
	slicePtr := (*[]interface{})(p)

	return *slicePtr
}

// apply writing
func (c *CopyOnWriteSlice) Apply(newSlice *[]interface{}) {
	atomic.StorePointer(&c.ptr_snapshot, unsafe.Pointer(newSlice))
}

// remove closed clients, clean?
// this func should run at intervals like 30 seconds
func (c *CopyOnWriteSlice) copyOnClean() {
	oldSlice := c.Snapshot()
	newSlice := make([]interface{}, 0, cap(oldSlice))

	for i := 0; i < len(oldSlice); i++ {
		if c.filter(oldSlice[i]) {
			newSlice = append(newSlice, oldSlice[i])
		}
	}

	c.Apply(&newSlice)
}

// add new element
// it looks like not a `copy-on-write` mode, but actually, partly it still does , due to the mechanism of `append`
func (c *CopyOnWriteSlice) copyOnAdd(e interface{}) {
	oldSlice := c.Snapshot()
	newSlice := append(oldSlice, e)

	c.Apply(&newSlice)
}

// write goroutine
func (c *CopyOnWriteSlice) loopWrite(clean_interval time.Duration) {
	clean_ticker := time.NewTicker(clean_interval)
	defer clean_ticker.Stop()

	for {
		select {
		case _, ok := <-c.closed:
			if !ok {
				log.Println("---------loopWrite closed--------")
				return
			}
		case e, ok := <-c.addChan:
			if ok {
				c.copyOnAdd(e)
				log.Println("< copyOnAdd actioned in loopWrite >", e)
			}
		case <-clean_ticker.C:
			select {
			case _, ok := <-c.cleanMark:
				if ok {
					c.copyOnClean()
					log.Println("< ------------Clean actioned in loopWrite >")
				}
			default:
				log.Println("+++ skip clean")
			}
		}
	}
	log.Println("loopWrite ended")
}

func (c *CopyOnWriteSlice) Add(e interface{}) {
	c.addChan <- e
}

func (c *CopyOnWriteSlice) MarkClean() {
	select {
	case c.cleanMark <- true:
	default:
	}
}
