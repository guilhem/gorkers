package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"

	"github.com/guilhem/gorkers"
)

var (
	count = make(map[string]int)
	mut   = sync.RWMutex{}
)

func main() {
	ctx := context.Background()

	workerOne := gorkers.NewRunner(ctx, NewWorkerOne().Work, 1000, 1000)
	workerTwo := gorkers.NewRunner(ctx, NewWorkerTwo().Work, 1000, 1000).InFrom(workerOne)
	if err := workerOne.Start(); err != nil {
		log.Panicf("workerOne start: %v", err)
	}
	if err := workerTwo.Start(); err != nil {
		log.Panicf("workerTwo start: %v", err)
	}

	for i := 0; i < 100000; i++ {
		workerOne.Send(rand.Intn(100))
	}
	workerOne.Wait().Stop()
	workerTwo.Wait().Stop()

	fmt.Println("worker_one", count["worker_one"])
	fmt.Println("worker_two", count["worker_two"])
	fmt.Println("finished")
}

type WorkerOne struct {
}
type WorkerTwo struct {
}

func NewWorkerOne() *WorkerOne {
	return &WorkerOne{}
}

func NewWorkerTwo() *WorkerTwo {
	return &WorkerTwo{}
}

func (wo *WorkerOne) Work(_ context.Context, in int, out chan<- int) error {
	var workerOne = "worker_one"
	mut.Lock()
	if val, ok := count[workerOne]; ok {
		count[workerOne] = val + 1
	} else {
		count[workerOne] = 1
	}
	mut.Unlock()

	total := in * 2
	fmt.Println("worker1", fmt.Sprintf("%d * 2 = %d", in, total))
	out <- total
	return nil
}

func (wt *WorkerTwo) Work(_ context.Context, in int, out chan<- interface{}) error {
	var workerTwo = "worker_two"
	mut.Lock()
	if val, ok := count[workerTwo]; ok {
		count[workerTwo] = val + 1
	} else {
		count[workerTwo] = 1
	}
	mut.Unlock()

	totalFromWorkerOne := in
	fmt.Println("worker2", fmt.Sprintf("%d * 4 = %d", totalFromWorkerOne, totalFromWorkerOne*4))
	return nil
}
