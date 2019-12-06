package async

import (
	"errors"
	"fmt"
	"sync"
)

// Handler is a function type for specifing the business logic of the task
type Handler func([]interface{}) TaskResult

// TaskResolver is a callback function type for handling the task result
type TaskResolver func(TaskResult)

// TasksetResolver is a callback function type for handling the taskset result
type TasksetResolver func([]TaskResult)

// Task is a struct which supports an asynchronous job.
type Task struct {
	handler Handler
	params  []interface{}
}

// Taskset is a struct which supports multiple parallel asynchronous job.
type Taskset struct {
	tasks []*Task
}

// TaskResult is a struct which contains result of the task.
type TaskResult struct {
	Data interface{}
	Err  error
}

// New creates a new Task object.
func New(handler Handler, params []interface{}) *Task {
	return &Task{
		handler: handler,
		params:  params,
	}
}

// Run runs handler function of the task asynchronously.
func (task *Task) Run(resolver TaskResolver) {
	go func(t *Task) {
		defer func() {
			if r := recover(); r != nil {
				resolver(TaskResult{
					Data: nil,
					Err:  fmt.Errorf(fmt.Sprintf("%s", r))})
			}
		}()
		resolver(t.handler(task.params))
	}(task)
}

// Await runs handler function of the task synchronously.
func (task *Task) Await() (res TaskResult) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	task.Run(func(r TaskResult) {
		defer wg.Done()
		res = r
	})
	wg.Wait()
	return
}

// Add adds a new Task to the Taskset
func (taskset *Taskset) Add(task *Task) {
	taskset.tasks = append(taskset.tasks, task)
}

// Run runs all handler function of the Taskset parallel asynchronously.
func (taskset *Taskset) Run(resolver TasksetResolver) {
	results := make([]TaskResult, len(taskset.tasks))

	go func(ts *Taskset) {
		wg := sync.WaitGroup{}
		wg.Add(len(ts.tasks))
		for idx := range ts.tasks {
			go func(task *Task, tid int) {
				defer wg.Done()
				if task == nil {
					results[tid] = TaskResult{
						Data: nil,
						Err:  errors.New("task is null")}
				}
				results[tid] = task.Await()
			}(taskset.tasks[idx], idx)
		}
		wg.Wait()
		resolver(results)
	}(taskset)

	return
}

// AwaitAll runs all handler function of the Taskset parallel synchronously.
func (taskset *Taskset) AwaitAll() (results []TaskResult) {
	results = make([]TaskResult, len(taskset.tasks))
	wg := sync.WaitGroup{}
	wg.Add(len(taskset.tasks))
	for idx := range taskset.tasks {
		go func(task *Task, tid int) {
			defer wg.Done()
			if task == nil {
				results[tid] = TaskResult{
					Data: nil,
					Err:  errors.New("task is null")}
			}
			results[tid] = task.Await()
		}(taskset.tasks[idx], idx)
	}
	wg.Wait()
	return
}
