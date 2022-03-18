package main

import (
	"context"
	"fmt"
	"time"

	"github.com/guilhem/gorkers"
)

func main() {
	ctx := context.Background()

	timeoutWorker := gorkers.NewRunner(ctx, work, 10, 10).SetWorkerTimeout(100 * time.Millisecond)
	timeoutWorker.AfterFunc(gorkers.StopWhenError[string])
	err := timeoutWorker.Start()
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < 1000000; i++ {
		timeoutWorker.Send("hello")
	}

	timeoutWorker.Wait().Stop()
}

func work(ctx context.Context, in string, out chan<- interface{}) error {
	fmt.Println(in)
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("overslept")
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
