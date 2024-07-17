package main

import "fmt"

type semaphore int	// 定义信号量类型

func P(S *semaphore){
	for *S<=0 {
	}
	*S-=1
}

func V(S *semaphore){
	*S+=1
}

var count int =0	// 用于记录当前的读者数目
var mutex semaphore = 1	// 用于保护更新count变量时的互斥
var rw semaphore = 1	// 用于保护读者与写者互斥地访问文件



func writer(rw *semaphore) {
	for {
		P(rw)
		fmt.Println("writing")
		V(rw)
	}
}

func reader(rw *semaphore, mutex *semaphore, count *int) {
	for{
		P(mutex)
		if *count==0 {
			P(rw)	
		}
		*count++
		V(mutex)
		fmt.Println("reading")
		P(mutex)
		*count--
		if *count==0 {
			V(rw)	// 允许写进程写
		}
		V(mutex)
	}
}
func main(){

	go reader(&rw, &mutex, &count)
	go writer(&rw)

}