package main

import (
	"fmt"
	"time"
)

func runCmd(input chan string, output chan string, command string)  {
	// read from channel via async
		go func() {
			for msg := range input {
				println("async function", msg)
			}
		}()

		for {

			println("goroutine", command)
			time.Sleep(time.Duration(2000) * time.Millisecond)
			output <- "from function"
		}
}

func controlledReader(name string, output <-chan string, control <-chan struct{}, next chan<- struct{}) {
	<-control // wait for permission
	counter := 0
	fmt.Println(name, "has control of output channel")

	for msg := range output {
		fmt.Println(name, "received:", msg)
		time.Sleep(500 * time.Millisecond) // simulate work
		counter += 1
		if counter >= 3 {
			fmt.Println(name, "exiting")
	next <- struct{}{} // pass control to next reader

		}
	}

	}

func main() {
	hello := "test"
	counter := 0
	output := make(chan string)
	input := make(chan string)
	
	// start async function
	go runCmd(input, output, hello)

	for {
		// // read from async function and display
		// go func() {
		// 	for msg := range output {
		// 		println(msg)
		// 	}
		// }()
		// // it is possible to spawn a second go routine that reads the same channel and sometimes takes the message instead of the first routine
		// go func() {
		// 	for msg := range output {
		// 		println("routine2", msg)
		// 	}
		// }()

		// control channels
		controlA := make(chan struct{}, 1)
		controlB := make(chan struct{}, 1)

		// reader A
		go controlledReader("ReaderA", output, controlA, controlB)

		// reader B
		go controlledReader("ReaderB", output, controlB, make(chan struct{})) // no handoff after B

		// give control to ReaderA initially
		controlA <- struct{}{}

		hello = hello + fmt.Sprintf("%d", counter)
		println(hello)
		input <- "1 via channel " + fmt.Sprintf("%d", counter)
		time.Sleep(time.Duration(1000) * time.Millisecond)
		input <- "2 via channel " + fmt.Sprintf("%d", counter)

		time.Sleep(time.Duration(7000) * time.Millisecond)
		counter += 1
	}
}
