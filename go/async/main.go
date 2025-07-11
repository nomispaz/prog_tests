package main

import (
	"fmt"
	"time"
)

func runCmd(input chan string, output chan string, command string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic while sending:", r)
		}
	}()
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

func main() {
	hello := "test"
	counter := 0
	output := make(chan string)
	input := make(chan string)

	// start async function
	go runCmd(input, output, hello)

	// read from async function and display
	go func() {
		for msg := range output {
			println(msg)
		}
	}()
	// it is possible to spawn a second go routine that reads the same channel and sometimes takes the message instead of the first routine
	go func() {
		for msg := range output {
			println("routine2", msg)
		}
	}()

	for {

		hello = hello + fmt.Sprintf("%d", counter)
		println(hello)
		input <- "1 via channel " + fmt.Sprintf("%d", counter)
		time.Sleep(time.Duration(1000) * time.Millisecond)
		input <- "2 via channel " + fmt.Sprintf("%d", counter)

		// it is possible to close the channel while a goroutine is still reading. With this, the program panics. But: not the goroutine panics but the sender in this loop
		close(output)
		time.Sleep(time.Duration(7000) * time.Millisecond)
		counter += 1
	}
}
