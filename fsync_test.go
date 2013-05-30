package fsync

import (
  // "log"
	"os"
	"strconv"
	"syscall"
	"testing"
	"time"
)

var (
	data     = []byte(`Ground control to major Tom. This is major Tom to ground control, I'm stepping through the door.`)
	flen     = int64(1000000)
	startPos = int64(100)
	nPerFile = 1000
)

func BenchmarkNosync(b *testing.B) {
	b.StopTimer()
	fname := "nosync." + strconv.FormatInt(time.Now().Unix(), 10)
	//log.Println("Create file ", fname)
	f, err := os.Create(fname)
	if err != nil {
		b.Fatal(err)
	}

	// if err := f.Truncate(flen); err != nil {
	// 	b.Fatal(err)
	// }

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		f.Seek(startPos, 0)
		for i := 0; i < nPerFile; i++ {
			f.Write(data)
		}
	}

}

func BenchmarkNosyncInitSize(b *testing.B) {
	b.StopTimer()
	fname := "nosync.sizeinit." + strconv.FormatInt(time.Now().Unix(), 10)
	//log.Println("Create file ", fname)
	f, err := os.Create(fname)
	if err != nil {
		b.Fatal(err)
	}

	if err := f.Truncate(flen); err != nil {
		b.Fatal(err)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		f.Seek(startPos, 0)
		for i := 0; i < nPerFile; i++ {
			f.Write(data)
		}
	}
}

func BenchmarkFsync(b *testing.B) {
	b.StopTimer()
	fname := "fsync." + strconv.FormatInt(time.Now().Unix(), 10)
	//log.Println("Create file ", fname)
	f, err := os.Create(fname)
	if err != nil {
		b.Fatal(err)
	}

	// if err := f.Truncate(flen); err != nil {
	// 	b.Fatal(err)
	// }

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		f.Seek(startPos, 0)
		for i := 0; i < nPerFile; i++ {
			f.Write(data)
			f.Sync()
		}
	}

}

func BenchmarkFsyncInitSize(b *testing.B) {
	b.StopTimer()
	fname := "fsync.sizeinit." + strconv.FormatInt(time.Now().Unix(), 10)
	//log.Println("Create file ", fname)
	f, err := os.Create(fname)
	if err != nil {
		b.Fatal(err)
	}

	if err := f.Truncate(flen); err != nil {
		b.Fatal(err)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		f.Seek(startPos, 0)
		for i := 0; i < nPerFile; i++ {
			f.Write(data)
			f.Sync()
		}
	}

}

func BenchmarkFdatasync(b *testing.B) {
	b.StopTimer()
	fname := "fdatasync." + strconv.FormatInt(time.Now().Unix(), 10)
	//log.Println("Create file ", fname)
	f, err := os.Create(fname)
	if err != nil {
		b.Fatal(err)
	}

	if err := f.Truncate(flen); err != nil {
		b.Fatal(err)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		f.Seek(startPos, 0)
		for i := 0; i < nPerFile; i++ {
			f.Write(data)
			Datasync(f)
		}
	}

}

func Datasync(f *os.File) (err error) {
	if f == nil {
		return syscall.EINVAL
	}
  // below only runs in *inux
	// if e := syscall.Fdatasync(int(f.Fd())); e != nil {
	// 	return os.NewSyscallError("fdatasync", e)
	// }
	return nil
}
