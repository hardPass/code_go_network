package main

import (
  "bytes"
	"container/list"
	"fmt"
	"github.com/hardPass/binary"
	"log"
	"net"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	Timeout  = time.Second * 10
	LogQueue = make(chan interface{}, 100)
	Exitlog  = make(chan true)
	Count    uint32
)

// BizlogicQueue <- &IoMsg{c, data}
// type IoMsg struct {
// 	C   *Client
// 	Msg []byte
// }

type Client struct {
	Name     string
	Conn     net.Conn
	Quit     chan bool
	Outgoing chan []byte
	BizLogic func(*Client, []byte) error
}

func (c *Client) Close() {
	c.Quit <- true
	//TODO unfinished
}

func CommonBizLogic(c *Client, data []byte) error {
	received := string(data)
	Log(received)
	c.Outgoing <- []byte("echo:" + received)
}

func (c *Client) LoopReceive() error {
	for {
		c.Conn.SetReadDeadline(time.Now().Add(Timeout))
		data, err := binary.Receive(c.Conn)
		if err != nil {
			return err
		}
		c.BizLogic(c, data)
	}
}

func (c *Client) LoopSend() error {
	for {
		select {
		case v := <-c.Outgoing:
			_, err := binary.Send(c, v)
			if err != nil {
				return err
			}
		case <-c.Quit:
			c.Conn.Close()
			break
		}
	}
	return nil
}

func Log(v interface{}) {
	LogQueue <- v
}

func LogServ() {
	file, err := os.OpenFile("tcp.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(-1)
	}
	defer file.Close()
	logger := log.New(file, "\r\n", log.Ldate|log.Ltime|log.Llongfile)

	go func() {
		for {
			select {
			case v := <-LogQueue:
				logger.Println(v)
			case <-Exitlog:
				logger.Println("----- Exitlog --------")
				return
			}
		}
	}()
}

func ClientHandler(conn net.Conn) {
	id := atomic.AddUint32(&Count, 1)
	c := &Client{
		Name:     strconv.FormatUint(id, 10),
		Conn:     conn,
		Quit:     make(chan bool),
		Outgoing: make(chan []byte, 2),
		BizLogic: CommonBizLogic,
	}
	go LoopSend(c)
	go LoopReceive(c)
}

func main() {
	LogServ()
	Log("----- Server started. ----- ")
	netListen, err := net.Listen("tcp", ":9988")
	if err != nil {
		Log(err)
		os.Exit(-1)
	}
	defer netListen.Close()

	for {
		conn, err := netListen.Accept()
		if err != nil {
			Log("Client error: ", err)
		} else {
			go IoHandler(conn)
		}
	}

}
