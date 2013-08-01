package main

import (
  "fmt"
)

type Bitmap struct {
	Shift uint32
	Mask  uint32

	data []uint32
}

func NewBitmap() *Bitmap {
	bm := &Bitmap{
		Shift: 5,
		Mask:  0x1F,
	}
	bm.data = make([]uint32, 16)
	// bm.data = make([]uint32, 1<<10)

	return bm
}

func (bm *Bitmap) Put(momoid uint32) {
	idx := momoid >> bm.Shift

	if oldSize := len(bm.data); int(idx) >= oldSize {
		tmp := make([]uint32, oldSize<<1)
		for i := 0; i < oldSize; i++ {
			tmp[i] = bm.data[i]
		}
		bm.data = tmp
	}

	bm.data[idx] |= (1 << (momoid & bm.Mask))
}

func (bm *Bitmap) Contains(momoid uint32) bool {
	idx := momoid >> bm.Shift
	if int(idx) >= len(bm.data) {
		return false
	}

	n := uint32(1 << (momoid & bm.Mask))

	return bm.data[idx]&n == n
}

func main() {
	bm := NewBitmap()
	bm.Put(0)
	bm.Put(11)
	bm.Put(21)
	bm.Put(12)
	bm.Put(123)
	fmt.Println(bm.Contains(0))
	fmt.Println(bm.Contains(11))
	fmt.Println(bm.Contains(21))
	fmt.Println(bm.Contains(12))
	fmt.Println(bm.Contains(123))
	fmt.Println(bm.Contains(511))
	fmt.Println(bm.Contains(512))
	fmt.Println(bm.Contains(513))
}
