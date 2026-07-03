package main

import (
	"fmt"
	"sync"
)

func main() {
	count := 0
	var wg sync.WaitGroup
	var mu sync.Mutex

	for range 1000 {
		wg.Add(1)
		wg.Go(func() {
			defer wg.Done()
			mu.Lock()
			count++
			mu.Unlock()
		})
	}
	wg.Wait()
	fmt.Println("final count:", count)
}
