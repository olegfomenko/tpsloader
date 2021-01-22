package pool

import (
	"context"
	"log"
	"sync"
	"time"
)

var (
	index      = workerIndex{index: 0}
	Timestamps []TransactionTimestamp
	Successful = 0
	Failed     = 0
)

type TransactionTimestamp struct {
	Start  time.Time
	Finish time.Time
	Status bool
}

type workerIndex struct {
	mu    sync.Mutex
	index int
}

func getNextId() int {
	index.mu.Lock()
	defer index.mu.Unlock()

	index.index++
	return index.index
}

type Task interface {
	Run(ch chan struct{})
}

type Pool interface {
	Submit(ctx context.Context, task ...Task)
	startTask(ctx context.Context, task Task)
}

type PoolImpl struct{}

func (pool *PoolImpl) startTask(ctx context.Context, task Task) {
	taskId := getNextId()
	log.Println("Starting task #", taskId)

	ch := make(chan struct{}, 1)
	go task.Run(ch)

	select {
	case <-ctx.Done():
		log.Println("Context finished. Finishing task #", taskId)
		ch <- struct{}{}
	case <-ch:
		log.Println("Task #", taskId, "finished himself")
	}

	log.Println("Task #", taskId, "finished!")
}

func (pool *PoolImpl) Submit(ctx context.Context, tasks ...Task) {
	for _, task := range tasks {
		go pool.startTask(ctx, task)
	}
}
