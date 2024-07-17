package main

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	value int
	cond  *sync.Cond	// 条件变量
}

func NewSemaphore(initial int) *Semaphore {
	return &Semaphore{
		value: initial,
		cond:  sync.NewCond(&sync.Mutex{}),
	}
}

func (s *Semaphore) P() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()
	for s.value <= 0 {
		s.cond.Wait()
	}
	s.value--
}

func (s *Semaphore) V() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()
	s.value++
	s.cond.Signal() // 唤醒一个等待在条件变量上的协程
}

var count int = 0           // 用于记录当前的读者数目
var mutex = NewSemaphore(1) // 用于保护更新count变量时的互斥
var rw = NewSemaphore(1)    // 用于保护读者与写者互斥地访问文件

func writer(rw *Semaphore) {
	for {
		rw.P()
		fmt.Println("writing")
		rw.V()
		time.Sleep(1 * time.Second) // 防止死循环导致输出过快
	}
}

func reader(rw *Semaphore, mutex *Semaphore, count *int) {
	for {
		mutex.P()
		if *count == 0 {
			rw.P()
		}
		*count++
		mutex.V()

		fmt.Println("reading")

		mutex.P()
		*count--
		if *count == 0 {
			rw.V()
		}
		mutex.V()

		time.Sleep(1 * time.Second)
	}
}

func main() {
	go reader(rw, mutex, &count)
	go writer(rw)

	// 让主协程等待，以便观察输出
	select {}	// 阻塞主协程
}
