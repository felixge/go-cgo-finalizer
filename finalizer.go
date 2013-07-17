package main

//#include <stdlib.h>
import "C"

import (
	"fmt"
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

const (
	loopSize = 10000
	sampleSize = 1000
)
var (
	rUsage syscall.Rusage
	allocs int
	frees int
	lock sync.Mutex
)

func main() {
	fmt.Printf("sample\trss_mb\tallocs\tfrees\n")
	for sample := 1; sample < sampleSize; sample++ {
		for i := 0; i <= loopSize; i++ {
			a := NewAllocator(10*1024)
			_ = a
		}

		err := syscall.Getrusage(syscall.RUSAGE_SELF, &rUsage)
		if err != nil {
			panic(err)
		}

		lock.Lock()
		fmt.Printf("%d\t%d\t%d\t%d\n", sample, rUsage.Maxrss/1024/1024, allocs, frees)
		lock.Unlock()
	}
}

func NewAllocator(size int) *Allocator {
	a := &Allocator{}
	a.alloc(size)
	runtime.SetFinalizer(a, free)
	return a
}

type Allocator struct{
	memory *C.char
}

func (a *Allocator) alloc(size int) {
	str := string(make([]byte, size))
	a.memory = C.CString(str)
	lock.Lock()
	defer lock.Unlock()
	allocs++
}

func free(a *Allocator) {
	C.free(unsafe.Pointer(a.memory))
	lock.Lock()
	defer lock.Unlock()
	frees++
}
