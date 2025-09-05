package main

import (
	"fmt"
	"sync"
)

var (
	counter int
	mutex   sync.RWMutex
	wg      sync.WaitGroup
)

func increment() {
	// mutex.Lock() // only single coroutine can use be here at a time

	// counter++
	// fmt.Println("Counter", counter)
	// mutex.Unlock()

	defer wg.Done()
	mutex.Lock()
	counter++
	fmt.Println("Counter", counter)
	mutex.Unlock()
}

func main() {
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go increment()
	}

	// giving time to finish the increase func
	// time.Sleep(2 * time.Second)

	wg.Wait()

	fmt.Println("Final Counter:", counter)
}
