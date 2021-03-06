package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/guilhem/gorkers"
)

func main() {
	ctx := context.Background()
	workerOne := gorkers.NewRunner(ctx, NewWorkerOne(2).Work, 100, 100)
	workerTwo := gorkers.NewRunner(ctx, NewWorkerTwo(4).Work, 100, 100).InFrom(workerOne)
	if err := workerOne.Start(); err != nil {
		fmt.Println(err)
	}

	if err := workerTwo.Start(); err != nil {
		log.Panic(err)
	}

	for i := 0; i < 15; i++ {
		workerOne.Send(rand.Intn(100))
	}

	workerOne.Wait().Stop()
	workerTwo.Wait().Stop()
}

type WorkerOne struct {
	amountToMultiply int
}
type WorkerTwo struct {
	amountToMultiply int
}

func NewWorkerOne(amountToMultiply int) *WorkerOne {
	return &WorkerOne{
		amountToMultiply: amountToMultiply,
	}
}

func NewWorkerTwo(amountToMultiply int) *WorkerTwo {
	return &WorkerTwo{
		amountToMultiply,
	}
}

func (wo *WorkerOne) Work(ctx context.Context, in int, out chan<- int) error {
	total := in * wo.amountToMultiply
	fmt.Println("worker1", fmt.Sprintf("%d * %d = %d", in, wo.amountToMultiply, total))
	out <- total
	return nil
}

func (wt *WorkerTwo) Work(ctx context.Context, in int, out chan<- interface{}) error {
	totalFromWorkerOne := in
	fmt.Println("worker2", fmt.Sprintf("%d * %d = %d", totalFromWorkerOne, wt.amountToMultiply, totalFromWorkerOne*wt.amountToMultiply))
	return nil
}
