package p2p

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
)

type OverflowPolicy int

const (
	Drop OverflowPolicy = iota
	Block
	Error
)

type Worker[T any] struct {
	queue   chan T
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	handler func(T)
	policy  OverflowPolicy
}

func NewWorker[T any](
	bufferSize int,
	ctx context.Context,
	policy OverflowPolicy,
	handler func(T),
) *Worker[T] {
	ctx, cancel := context.WithCancel(ctx)
	return &Worker[T]{
		queue:   make(chan T, bufferSize),
		ctx:     ctx,
		cancel:  cancel,
		handler: handler,
		policy:  policy,
	}
}

func (w *Worker[T]) Start(concurrency int) {
	for i := range concurrency {
		w.wg.Add(1)
		go func(id int) {
			defer w.wg.Done()
			for {
				select {
				case <-w.ctx.Done():
					log.Infof("Worker-%d shutting down...", id)
					return
				case job := <-w.queue:
					if w.handler != nil {
						w.handler(job)
					}
				}
			}
		}(i)
	}
}

func (w *Worker[T]) Push(job T) {
	switch w.policy {
	case Drop:
		select {
		case w.queue <- job:
		default:
			<-w.queue
			log.Warn("Worker queue full, drop job")
		}
	case Block:
		w.queue <- job
	case Error:
		select {
		case w.queue <- job:
		default:
			log.Error("Worker queue full")
		}
	}
}

func (w *Worker[T]) Shutdown() {
	w.cancel()
	w.wg.Wait()
	close(w.queue)
}
