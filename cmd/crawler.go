package main

import (
	"fmt"
	uuid2 "github.com/google/uuid"
	"slot-crawler/internal/crawler"
	"sync"
	"time"
)

func main() {

	//database.Initialize()

	uuid := uuid2.New().String()
	slot, err := crawler.Initialize("vswaysdogs", uuid)
	if err != nil {
		fmt.Printf("error!!! : %s\n", err.Error())
	}

	var wg sync.WaitGroup
	counter := make(chan int)
	wg.Add(1)

	go func() {
		err := slot.StartCrawling(1000, time.Millisecond*300, counter)
		if err != nil {
			fmt.Println(err.Error())
		}
		wg.Done()
	}()

	go func() {
		for val := range counter {
			fmt.Println("Counter:", val)
		}
	}()

	wg.Wait()
	close(counter)
	fmt.Println("finished!!!!!!")

}
