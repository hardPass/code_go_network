package main

import (
  "bytes"
	"container/list"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	line_tpl = ` - - [26/Mar/2012:19:36:50 +0800] "GET /bbs/ck.php? HTTP/1.1" 200 296 "http://www.baidu.com/bbs/" "Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 5.1; Trident/4.0; KB974489)"`
)

var (
	minIP         = 1 << 24
	r             = rand.New(rand.NewSource(time.Now().Unix()))
	toDisk        = make(chan []byte, 10)
	piece         = 1 << 16
	round         = 1 << 29 // change rand  by seed every round
	rounds        = 0
	repeatListMax = 1 << 15
	maxSize       = 1 << 30 * 100
	total         = 0 // total log file's size actually
	allLines      = 0
	repeated      = 0
	torepeated    = 0

	repeat *list.List
	start  int64 // start time
)

func init() {
	repeat = list.New()
	// fmt.Println(minIP)
	// fmt.Println(maxSize)
}

func chance() bool {
	if r.Int()&1 == 0 {
		return true
	}

	return false
}

func ip4(v uint32) (string, string, string, string) {
	b := make([]byte, 4)
	b[0] = byte(v >> 24)
	b[1] = byte(v >> 16)
	b[2] = byte(v >> 8)
	b[3] = byte(v)

	return strconv.FormatInt(int64(b[0]), 10),
		strconv.FormatInt(int64(b[1]), 10),
		strconv.FormatInt(int64(b[2]), 10),
		strconv.FormatInt(int64(b[3]), 10)
}

func nextIP() uint32 {
	if chance() {
		e := repeat.Front()
		if e != nil {
			// n := uint32(repeat.Remove(e))
			n, _ := repeat.Remove(e).(uint32)
			repeated++

			if chance() && repeat.Len() < repeatListMax+1000 {
				torepeated++
				repeat.PushBack(n)
			}

			return n
		}
	}

	for {
		n := r.Uint32()
		if int(n) > minIP {
			if repeat.Len() < repeatListMax {
				torepeated++
				repeat.PushBack(n)
			}

			return n
		}
	}
}

func logPiece() ([]byte, int) {
	buf := bytes.NewBuffer(make([]byte, 0, piece+300))
	for {
		allLines++
		ip0, ip1, ip2, ip3 := ip4(nextIP())
		buf.WriteString(ip0)
		buf.WriteByte('.')
		buf.WriteString(ip1)
		buf.WriteByte('.')
		buf.WriteString(ip2)
		buf.WriteByte('.')
		buf.WriteString(ip3)
		buf.WriteString(line_tpl)
		buf.WriteByte('\n')
		if buf.Len() >= piece {
			return buf.Bytes(), buf.Len()
		}
	}
}

func logging() {
	for {
		b, len := logPiece()
		// fmt.Println("log a piece:", len)
		toDisk <- b
		total += len
		if total >= maxSize {
			toDisk <- make([]byte, 0)
			return
		}

		if total > (rounds * round) {
			rounds++

			r = rand.New(rand.NewSource(time.Now().Unix()))
			fmt.Printf("proceed %d MB , all spend: %d ms \n", total/(1<<20), (time.Now().UnixNano()-start)/(1000*1000))
		}
	}
	fmt.Println("logging done")

}

func toFs(fi *os.File) {
	for {
		select {
		case b, _ := <-toDisk:
			// fmt.Println("to fs:", len(b))
			if len(b) == 0 {
				return
			}

			fi.Write(b)
		}
	}
	fmt.Println("toFs done")
}

func main() {
	fi, err := os.OpenFile("./100g.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0420)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()
	fi.Seek(0, 0)

	runtime.GOMAXPROCS(3)

	fmt.Println("Begin...")
	start = time.Now().UnixNano()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		logging()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		toFs(fi)
	}()

	wg.Wait()
	end := time.Now().UnixNano()

	fmt.Printf("Done in %d ms!\n", (end-start)/(1000*1000))
	fmt.Printf("allLines %d\n", allLines)
	fmt.Printf("total %d\n", total)
	fmt.Printf("repeated %d\n", repeated)
	fmt.Printf("torepeated %d\n", torepeated)

}
