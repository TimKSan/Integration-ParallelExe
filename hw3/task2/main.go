package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		i := 2
		for {
			select {
			case <-done:
				fmt.Println("\nПоступил сигнал остановки")
				return
			case <-time.After(1 * time.Second):
				if i%2 == 0 {
					fmt.Println("Квадрат натурального числа:", i*i)
				}
				i++
			}
		}
	}()

	wg.Wait()
	fmt.Println("Выхожу из программы")
}
