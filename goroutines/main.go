/*Concurrency and Parallelism*/
/*Concurrency: getting items from a shopping list in a store
single task at a time*/
/*Parallelism: now you and your friends shop with you in a store
buying different items from the list (built on top of concurrency)
multiple tasks at the same time*/
package main

import (
	"fmt"
	"math/rand"
	"time"
)

var cache = map[int]Book{}
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	for i := 0; i < 10; i++ {
		id := rnd.Intn(10) + 1 //random int between 0 and 9
		if b, ok := queryCache(id); ok {
			fmt.Println("form cache")
			fmt.Println(b)
			//else if or `continue` here
		} else if b, ok := queryDatabase(id); ok {
			fmt.Println("form database")
			fmt.Println(b)
			//else of `continue` here
		} else {
			fmt.Printf("Book not found with id: '%v'", id)
			time.Sleep(150 * time.Microsecond)
		}
	}

}

func queryCache(id int) (Book, bool) {
	b, ok := cache[id]
	return b, ok
}

func queryDatabase(id int) (Book, bool) {
	time.Sleep(100 * time.Millisecond)
	//assuming the that the books might not be sorted
	for _, b := range books {
		if b.ID == id {
			cache[id] = b
			return b, true
		}
	}

	return Book{}, false
}
