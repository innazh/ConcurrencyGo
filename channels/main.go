package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	ch := make(chan int)

	wg.Add(2)

	go func(ch <-chan int, wg *sync.WaitGroup) {
		for msg := range ch {
			//receiving a msg from the channel
			fmt.Println(msg)
		} //the for-loop statement is gonna iterate through the channel's msgs until it closes, looking for a close() statement to stop
		wg.Done()
	}(ch, wg)

	go func(ch chan<- int, wg *sync.WaitGroup) {
		for i := 0; i < 10; i++ {
			//sending 42 into the channel
			ch <- i
		}
		close(ch)
		wg.Done()
	}(ch, wg)

	wg.Wait()
}

/*Channel has a blocking inteface and blocks the execution until someone is available to either receive or send a msg (depending on what you do)
That's why goroutines are neccessary, one of them can sleep while the other one works and is able to send/receive a msg
If the number of senders and receivers doesn't match on an unbuffered channel, you get a deadlock condition*/
/*Channel types:*/
/*ch := make(chan int) -- always bidirectional by default*/
/*func f(ch chan int){} - bidirectional channel, can send and receive msgs*/
/*func f(ch chan<- int){} - send-only channel, can only send msgs*/
/*func f(ch <-chan int){} - receive-only channel, can only receive msgs*/

/*Can only close channel on the sending side, not on the receiving one. If you try to send a msg into the closed channel, the program will panic, if you try to receive it
from a closed channel, you'll get 0*/

/*if msg,ok := <-ch; ok {} , ok is gonna be true if the channel's still open, false otherwise*/
