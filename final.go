package main

import (
	"fmt"
	"sync"
)

func partial_factorial(n_prev int64, n int64, c chan int64, wg *sync.WaitGroup) {
	var partial_product int64
	partial_product = 1
	for i := n_prev; i <= n; i++ {
		if i > 1 {
			partial_product *= int64(i)
		}
	}
	c <- partial_product
	wg.Done()
}

func factorial(n int64, prev_map *map[int64]int64, prev_holder *int64, mut *sync.Mutex) int64 {
	prev_holder_exec := *prev_holder
	factorial_n := (*prev_map)[prev_holder_exec]

	var wg sync.WaitGroup
	var channel_size int64
	if n > prev_holder_exec {
		channel_size = (n - prev_holder_exec) / 5
	} else {
		channel_size = n / 5
	}
	channel_products := make(chan int64, int(channel_size)+1)

	if n > prev_holder_exec {
		wg.Add(int((n-prev_holder_exec)/5) + 1)
		for i := prev_holder_exec + 1; i <= n; i++ {
			if i%5 == 0 {
				go partial_factorial(i-4, i, channel_products, &wg)
			}
		}
	} else {
		wg.Add(int(n / 5))
		for i := int64(1); i <= n; i++ {
			if i%5 == 0 {
				go partial_factorial(i-4, i, channel_products, &wg)
			}
		}
	}

	wg.Wait()
	for v := range channel_products {
		factorial_n *= v
	}
	mut.Lock()
	if (n-prev_holder_exec > 10) && (n%10 == 0) {
		(*prev_map)[n/10] = factorial_n
	}
	mut.Unlock()

	return factorial_n
}

func partial_taylor(n int64, prev_map *map[int64]int64, prev_holder *int64, mut *sync.Mutex, c chan float64, wg *sync.WaitGroup) {
	var acc float64 = 0
	for i := n - 4; i <= n; i++ {
		acc += 1 / float64(factorial(i, prev_map, prev_holder, mut))
	}
	c <- acc
	wg.Done()
}

func calculate_taylor(n int64, prev_map *map[int64]int64, prev_holder *int64, mut *sync.Mutex) float64 {
	var taylor_accumulator float64
	var wg sync.WaitGroup

	taylor_accumulator = 1
	taylor_channel := make(chan float64, int(n%5))

	wg.Add(int(n / 5))

	for i := int64(0); i <= n; i++ {
		if i%5 == 0 {
			go partial_taylor(i, prev_map, prev_holder, mut, taylor_channel, &wg)
		}
	}

	wg.Wait()

	for v := range taylor_channel {
		taylor_accumulator += v
	}

	return taylor_accumulator
}

func main() {
	prev_map := make(map[int64]int64)
	var prev_holder int64
	var prev_mutex sync.Mutex

	fmt.Println(calculate_taylor(100, &prev_map, &prev_holder, &prev_mutex))
}
