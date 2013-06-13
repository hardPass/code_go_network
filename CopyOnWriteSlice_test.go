package concurrent

import (
  "bytes"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

var (
	list *CopyOnWriteSlice
)

type Client struct {
	id     int
	closed chan bool
}

// don't execute this twice
func (c *Client) Disconnect() {
	close(c.closed)
	list.MarkClean()
}

// Did you see the trick here?
func (c *Client) Closed() bool {
	select {
	case _, ok := <-c.closed:
		if !ok {
			return true
		}
	default:
	}

	return false
}

func (c *Client) String() string {
	// stat := "live"
	stat := "L"
	if c.Closed() {
		// stat = "dead"
		stat = "D"
	}
	return fmt.Sprintf("{%d, %v}", c.id, stat)
}

func TestCopyOnWriteSlice(t *testing.T) {

	filter := func(e interface{}) bool {
		c, _ := e.(*Client)
		return !c.Closed()
	}
	list = New(10000, filter, time.Millisecond*1000)

	var wg sync.WaitGroup

	// simulate client connect actions
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(time.Millisecond * 300)
		defer ticker.Stop()

		id := 0

		for {
			id++
			c := &Client{id: id, closed: make(chan bool)}
			list.Add(c)
			log.Println("new:", c)

			<-ticker.C
		}
	}()

	// simulate client disconnect actions
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(time.Millisecond * 300)
		defer ticker.Stop()

		for {
			clients := list.Snapshot()
			for i := 0; i < len(clients); i++ {
				c := clients[i].(*Client)
				if i%3 == 0 && !c.Closed() { // disconnect some clients
					c.Disconnect()
					log.Println("disconnect:", c)
				}
			}

			<-ticker.C
		}
	}()

	// moniter the state of all clients
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(time.Millisecond * 300)
		defer ticker.Stop()

		for {

			clients := list.Snapshot()
			// log.Println("all--", len(clients), cap(clients))s
			log.Println("all--", String(clients))

			<-ticker.C
		}
	}()

	wg.Wait()
	list.Close()

	time.Sleep(222)
	log.Println("-------------------exit")
}

func String(clients []interface{}) string {
	buffer := bytes.NewBuffer(make([]byte, 0, 1000))

	for i := 0; i < len(clients); i++ {
		buffer.WriteString(clients[i].(*Client).String())
		buffer.WriteString(", ")
	}

	return buffer.String()
}
