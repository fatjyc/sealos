package helper

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockTask struct {
	executed bool
	err      error
}

func (t *mockTask) Execute() error {
	t.executed = true
	return t.err
}

func TestNewTaskQueue(t *testing.T) {
	ctx := context.Background()
	tq := NewTaskQueue(ctx, 2, 10)

	assert.NotNil(t, tq)
	assert.Equal(t, 2, tq.workers)
	assert.Equal(t, 10, cap(tq.queue))
	assert.NotNil(t, tq.ctx)
	assert.NotNil(t, tq.cancel)
	assert.False(t, tq.started)
}

func TestTaskQueue_Start(t *testing.T) {
	tq := NewTaskQueue(context.Background(), 2, 10)

	// Start queue
	tq.Start()
	assert.True(t, tq.started)

	// Start again should not spawn more workers
	tq.Start()
	assert.True(t, tq.started)

	tq.Stop()
}

func TestTaskQueue_AddTask(t *testing.T) {
	ctx := context.Background()
	tq := NewTaskQueue(ctx, 2, 10)
	tq.Start()

	task := &mockTask{}
	tq.AddTask(task)

	// Give time for task to be processed
	time.Sleep(100 * time.Millisecond)

	assert.True(t, task.executed)

	tq.Stop()
}

func TestTaskQueue_Stop(t *testing.T) {
	ctx := context.Background()
	tq := NewTaskQueue(ctx, 2, 10)
	tq.Start()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		task := &mockTask{}
		tq.AddTask(task)
	}()

	// Stop queue
	tq.Stop()

	// Verify queue is stopped
	select {
	case _, ok := <-tq.queue:
		assert.False(t, ok)
	default:
	}

	wg.Wait()
}

func TestTaskQueue_Worker(t *testing.T) {
	ctx := context.Background()
	tq := NewTaskQueue(ctx, 2, 10)
	tq.Start()

	// Test successful task
	task1 := &mockTask{}
	tq.AddTask(task1)

	// Test failed task
	task2 := &mockTask{err: assert.AnError}
	tq.AddTask(task2)

	// Give time for tasks to be processed
	time.Sleep(100 * time.Millisecond)

	assert.True(t, task1.executed)
	assert.True(t, task2.executed)

	tq.Stop()
}

func TestTaskQueue_AddTask_AfterStop(t *testing.T) {
	ctx := context.Background()
	tq := NewTaskQueue(ctx, 2, 10)
	tq.Start()
	tq.Stop()

	task := &mockTask{}
	tq.AddTask(task)

	assert.False(t, task.executed)
}
