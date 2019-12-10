# Go Async
Simple async utility for supporting your multiple asynchronous tasks.

## Key Features
- Callback style handler
- Await multiple asynchronous tasks
- Combined results of multiple asynchronous tasks
- Panic safety goroutine

# Installation
    $ go get github.com/hisdream86/go-async

# Example
```go
package main

import (
	"fmt"
	"time"

	"github.com/hisdream86/go-async"
)

func handler1(params []interface{}) async.TaskResult {
	sec := params[0].(int64)

	time.Sleep(time.Duration(sec) * time.Second)

	return async.TaskResult{
		Data: "hello",
		Err:  nil,
	}
}

func handler2(params []interface{}) async.TaskResult {
	sec := params[0].(int64)
	str := params[1].(string)

	time.Sleep(time.Duration(sec) * time.Second)

	return async.TaskResult{
		Data: str,
		Err:  nil,
	}
}

func main() {
	task1 := async.New(handler1, []interface{}{int64(1)})
	task2 := async.New(handler2, []interface{}{int64(2), "world"})

	// Single async task
	// Print "hello" after 1 second
	task1.Run(func(result async.TaskResult) {
		fmt.Println(result.Data)
	})
	time.Sleep(2 * time.Second)

	var taskset async.Taskset
	taskset.Add(task1)
	taskset.Add(task2)

	// Multiple async tasks
	// Print "hello world" after 2 second
	taskset.Run(func(results []async.TaskResult) {
		fmt.Println((results[0].Data).(string) + " " + (results[1].Data).(string))
	})
	time.Sleep(3 * time.Second)

	// Await multiple async tasks
	// Print "hello world" after 2 second
	results := taskset.AwaitAll()
	fmt.Println((results[0].Data).(string) + " " + (results[1].Data).(string))
}
```