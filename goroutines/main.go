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

	/*WaitGroups helps us make the main function Wait for the goroutines to execute w/o needing to use time.Sleep()*/
	"sync"
	/*Mutex - mutual exclusion lock - helps us to manage shared memory*/
	"time"
)

var cache = map[int]Book{}
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	wg := &sync.WaitGroup{}
	m := &sync.RWMutex{}
	dbCh := make(chan Book)
	cacheCh := make(chan Book)

	for i := 0; i < 10; i++ {
		id := rnd.Intn(10) + 1 //random int between 0 and 9
		wg.Add(2)
		go func(id int, wg *sync.WaitGroup, m *sync.RWMutex, ch chan<- Book) {
			if b, ok := queryCache(id, m); ok {
				// fmt.Println("form cache")
				// fmt.Println(b)
				ch <- b
			}
			wg.Done()
		}(id, wg, m, cacheCh)

		go func(id int, wg *sync.WaitGroup, m *sync.RWMutex, ch chan<- Book) {
			if b, ok := queryDatabase(id, m); ok {
				// fmt.Println("form database")
				// fmt.Println(b)
				ch <- b
			}
			wg.Done()
		}(id, wg, m, dbCh)
		//The execution is complete once main finishes. This call is neccessary to give goroutines enough time to execute and return to the main function, otherwise we'd have no output
		//because main would finish its execution before the goroutines would
		// time.Sleep(150 * time.Millisecond)

		//create one goroutine per query to handle response
		go func(dbCh, cacheCh <-chan Book) {
			select {
			case b := <-dbCh:
				fmt.Println("from database")
				fmt.Println(b)
			case b := <-cacheCh:
				fmt.Println("from cache")
				fmt.Println(b)
				<-dbCh
			}
		}(dbCh, cacheCh)
		time.Sleep(150 * time.Millisecond) //that's for main to wait for this goroutine to execute
	}
	wg.Wait()
}

/*Use RWMutex only when you've got way more readers than writers, performance is worse than regular mutex*/
func queryCache(id int, m *sync.RWMutex) (Book, bool) {
	//allows multiple readers to read from the shared memory but when something is trying to write, then it clears out the readers and lets the writer come in
	m.RLock()
	b, ok := cache[id]
	m.RUnlock()
	return b, ok
}

func queryDatabase(id int, m *sync.RWMutex) (Book, bool) {
	time.Sleep(100 * time.Millisecond)
	//assuming the that the books might not be sorted
	for _, b := range books {
		if b.ID == id {
			//makes sure there are no readers (lets all current readers finish) and only 1 writer, accesses the memory, then unlocks
			m.Lock()
			cache[id] = b //should be a problem(race condition) because 2 goroutines are trying to access the same variable but it's not... go run --race .
			m.Unlock()
			return b, true
		}
	}

	return Book{}, false
}
