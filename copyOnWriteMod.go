package main

import (
  "fmt"
	// "runtime"
	"sync/atomic"
	"time"
	"unsafe"
)

// This program shows how a `copy-on-write` mode works.

var (
	seq              uint32
	clients_snapshot = unsafe.Pointer(&[]*Client{}) // *[]*Client   a pointer of a slice, just consider it as a snapshot
	newClientChan    = make(chan *Client, 10)
)

func nextSeq() uint32 {
	return atomic.AddUint32(&seq, 1)
}

// a very simple TCP client, which leaves out lots properties
// 一个简单的TCP客户端， 省略了一些东西，这边用 leave out 这个词合适不？
type Client struct {
	id    uint32
	valid chan bool
}

func newClient() *Client {
	return &Client{nextSeq(), make(chan bool)}
}

func (c *Client) Disconnect() {
	close(c.valid)
}

// Did you see the trick here?
func (c *Client) Closed() bool {
	select {
	case _, ok := <-c.valid:
		if !ok {
			return true
		}
	default:
	}

	return false
}

func (c *Client) String() string {
	stat := "live"
	if c.Closed() {
		stat = "dead"
	}
	return fmt.Sprintf("{%d, %v}", c.id, stat)
}

// only for read
func SnapshotClients() []*Client {
	p := atomic.LoadPointer(&clients_snapshot)
	slicePtr := (*[]*Client)(p)

	return *slicePtr
}

// remove closed clients, clean?
// this func should run at intervals like 30 seconds
func CopyOnRemove() {
	oldClients := SnapshotClients()
	newClients := []*Client{}

	// add all non-closed clients from snapshot, consider it as removing action

	for i := 0; i < len(oldClients); i++ {
		if !oldClients[i].Closed() {
			newClients = append(newClients, oldClients[i])
		}
	}

	atomic.StorePointer(&clients_snapshot, unsafe.Pointer(&newClients))
}

// add new client
// it looks like not a `copy-on-write` mode, but actually, partly it still does , due to the mechanism of `append`
func CopyOnAdd(c *Client) {
	oldClients := SnapshotClients()
	newClients := append(oldClients, c)

	atomic.StorePointer(&clients_snapshot, unsafe.Pointer(&newClients))
}

// write goroutine
func loopWrite() {
	rm_ticker := time.NewTicker(1e9 * 2) // for remove intervals
	defer rm_ticker.Stop()

	for {
		select {
		case c, ok := <-newClientChan:
			if ok {
				CopyOnAdd(c)
				fmt.Println("< add action >")
			}
		case <-rm_ticker.C:
			CopyOnRemove()
			fmt.Println("< clean action >")
		}
	}
}

func main() {

	go loopWrite()

	// simulate client connect actions
	go func() {
		ticker := time.NewTicker(1e8 * 3)
		defer ticker.Stop()

		for {
			c := newClient()
			newClientChan <- c
			fmt.Println("new:", c)

			<-ticker.C
		}
	}()

	// simulate client disconnect actions
	go func() {

		ticker := time.NewTicker(1e8 * 5)
		defer ticker.Stop()
		for {
			clients := SnapshotClients()
			for i := 0; i < len(clients); i++ {
				if i%3 == 0 && !clients[i].Closed() { // disconnect some clients
					clients[i].Disconnect()
					fmt.Println("disconnect:", clients[i])
				}
			}

			<-ticker.C
		}
	}()

	// moniter the state of all clients
	go func() {
		ticker := time.NewTicker(1e8 * 5)
		defer ticker.Stop()
		for {

			clients := SnapshotClients()
			fmt.Println("all--", clients)

			<-ticker.C
		}
	}()

	time.Sleep(1e11)
}
