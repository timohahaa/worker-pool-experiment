package manager

import (
	"context"

	"github.com/google/uuid"
	"github.com/timohahaa/sabbrum/task"
	"github.com/timohahaa/sabbrum/worker"
)

type Manager struct {
	workers            []*worker.Worker
	workerCount        int
	workerQueueSize    int
	completeTasksChan  chan *task.CompleteTask
	completeTasksMap   map[uuid.UUID]*task.CompleteTask
	roundRobbinCounter int
	ctx                context.Context
}

func New(ctx context.Context, workerCount, workerQueueSize int) *Manager {
	m := &Manager{
		workers:            make([]*worker.Worker, 0, workerCount),
		workerCount:        workerCount,
		workerQueueSize:    workerQueueSize,
		workerCancels:      make(map[*worker.Worker]context.CancelFunc),
		completeTasksChan:  make(chan *task.CompleteTask),
		roundRobbinCounter: 0,
		ctx:                ctx,
	}
	return m
}

func (m *Manager) InitWorkers() {
	for i := 0; i < m.workerCount; i++ {
	}
}

func (m *Manager) roundRobbinPutTask(t *task.Task) {
	m.roundRobbinCounter %= m.workerCount
	worker := m.workers[m.roundRobbinCounter]
	// this func is async waiting for a queue to have space in a buffered channel
	worker.AddTask(t)
}

func (m *Manager) getCompleteTaskByUUID(uuid uuid.UUID) *task.CompleteTask {
	t := m.completeTasksMap[uuid]
	delete(m.completeTasksMap, uuid)
	return t
}

func (m *Manager) Execute(t *task.Task) *task.CompleteTask {
	m.roundRobbinPutTask(t)

}
