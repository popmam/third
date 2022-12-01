package main

import (
	"fmt"
	//"sync"
	"time"
)

func printChannelData(c chan int) {
	now := time.Now()
	time.Sleep(5 * time.Second)
	fmt.Println(now.Format("15:04:05"), "in printeChannelData ", time.Now().Format("15:04:05"))
	c <- 25
}

func main() {
	now := time.Now()
	p := fmt.Println
	p("hello", now.Format("15:04:05"))
	c := make(chan int)
	go printChannelData(c)

	now = time.Now()
	p("after sleeping, main time is ", now.Format("15:04:05"), "but now it's ", time.Now().Format("15:04:05"), <-c)
	p("Done.")
}
