package main

import (
	"fmt"
)

type Bitmap struct {
	Shift uint32
	Mask  uint32
	Max   uint32

	Data []uint32
}

func NewBitmap() *Bitmap {
	bm := &Bitmap{
		Shift: 5,
		Mask:  0x1F,
	}
	// bm.Data = make([]uint32, 16)
	bm.Data = make([]uint32, 1<<29)

	return bm
}

func (bm *Bitmap) idx_of_ints(momoid uint32) uint32 {
	return momoid >> bm.Shift
}

func (bm *Bitmap) valueByOffset(momoid uint32) uint32 {
	return 1 << (momoid & bm.Mask)
}

func (bm *Bitmap) Put(momoid uint32) {
	if momoid > bm.Max {
		bm.Max = momoid
	}

	idx := bm.idx_of_ints(momoid)

	if oldSize := len(bm.Data); int(idx) >= oldSize {
		tmp := make([]uint32, idx+32)
		for i := 0; i < oldSize; i++ {
			tmp[i] = bm.Data[i]
		}
		bm.Data = tmp
	}

	v := bm.valueByOffset(momoid)
	bm.Data[idx] |= v
}

func (bm *Bitmap) Contains(momoid uint32) bool {
	idx := bm.idx_of_ints(momoid)
	if int(idx) >= len(bm.Data) {
		return false
	}

	v := bm.valueByOffset(momoid)

	return bm.Data[idx]&v == v
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
