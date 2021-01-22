package pool

import (
	"context"
	"log"
	"time"
)

var (
	Timestamps []time.Time
	Successful = 0
	Failed     = 0
)

type Task interface {
	Run(ctx context.Context, ready chan Task)
	GetName() string
}

type Pool interface {
	Submit(ctx context.Context, task ...Task)
	GetName() string
}

type poolImpl struct {
	name string
}

func (pool *poolImpl) Submit(ctx context.Context, tasks ...Task) {
	log.Println("Starting", pool.name)
	ready := make(chan Task, len(tasks))

	for _, task := range tasks {
		ready <- task
	}

	context, cancel := context.WithCancel(ctx)
	defer cancel()

	inProgress := true

	for inProgress {
		select {
		case task := <-ready:
			log.Println("Running task:", task.GetName())
			go task.Run(context, ready)
		case <-ctx.Done():
			log.Println("Main context closed. Closing pool:", pool.name)
			inProgress = false
		}
	}

	log.Println("Pool", pool.name, "finished")
}

func (pool *poolImpl) GetName() string {
	return pool.name
}

func GetPool(name string) Pool {
	return &poolImpl{
		name: name,
	}
}
