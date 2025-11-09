package main

import (
	"fmt"
	"time"
)

func worker(c chan string, i int) {
	for msg := range c {
		fmt.Println("worker", i, "received message", msg)
		time.Sleep(10 * time.Second)
	}

}

func main() {

	c := make(chan string)

	for i := 0; i < 3; i++ {
		go worker(c, i)
	}

	for i := 0; i < 100; i++ {
		c <- fmt.Sprintf("message %d", i)
	}

	time.Sleep(10 * time.Second)
	//close(c)

}
