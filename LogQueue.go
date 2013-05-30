package main

import (
  "bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"sync"
	"sync/atomic"
)

// type LogQueue interface {
// 	EnQueue([]byte) error
// 	DeQueue() ([]byte, uint32)
// 	Commit(n uint32) error
// 	Close() error
// 	Pause()
// 	Start()
// 	Len() int
// }

// chan() chan []byte // this is expected to be an *unbuffered* channel

type LogQueue struct {
	sync.Mutex

	MChan chan []byte 

	EnQueue([]byte) error
	DeQueue() ([]byte, uint32)
	Commit(n uint32) error
	Close() error
	Pause()
	Start()
	Len() int
}
