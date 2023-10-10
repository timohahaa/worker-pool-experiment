package worker

import (
	"context"
	"sync"

	"github.com/timohahaa/croc"
	"github.com/timohahaa/sabbrum/task"
)

// NOTE: why I decided to make cancelFunc available for a worker?
type Worker struct {
	client           *croc.CrocClient
	taskQueueSize    int
	inputQueue       chan *task.Task
	outputQueue      chan *task.CompleteTask
	tasksLeft        int
	ctx              context.Context
	cancelWorkerFunc context.CancelFunc
	mu               sync.RWMutex
}

// the way it is designed, you dont need to set input queue explicitly,
// instead you set output queue explicitly, and put tasks into input queue with AddTask() method
// worker is stopped using <-ctx.Done()
func New(ctx context.Context, queueSize int, output chan *task.CompleteTask) *Worker {
	childContext, cancel := context.WithCancel(ctx)
	w := &Worker{
		client:           croc.New(),
		taskQueueSize:    queueSize,
		inputQueue:       make(chan *task.Task, queueSize),
		outputQueue:      output,
		tasksLeft:        0,
		ctx:              childContext,
		cancelWorkerFunc: cancel,
		mu:               sync.RWMutex{},
	}
	return w
}

// this funtion does not need to make mutex calls, because all of the tasks are read in a queue maner
func (w *Worker) Run() {
	go func() {
		for {
			select {
			case t := <-w.inputQueue:
				// set proxy
				w.client.Proxy(t.Proxy)
				// make a request
				resp, _, err := w.client.Do(t.Request)
				// send complete task
				// set complete task's id to be the same as the task it was derived from
				completeTask := task.NewCompleteTask(t.Id, resp, err)
				w.outputQueue <- completeTask
				w.tasksLeft--
			case <-w.ctx.Done():
				// worker recieved a finish signal
				return
			}
		}
	}()
}

// use mutex here, because main Run() writes to w.tasksLeft
func (w *Worker) AddTask(t *task.Task) {
	go func() {
		w.inputQueue <- t
		w.mu.Lock()
		defer w.mu.Unlock()
		w.tasksLeft++
	}()
}

// use mutex here, because main Run() writes to w.tasksLeft
func (w *Worker) TasksLeft() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.tasksLeft

}
