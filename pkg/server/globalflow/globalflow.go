package globalflow

import (
	"sync"
	"time"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/typeopr"
)

type Task func(manager interfaces.Manager)

// GlobalFlow creates a flow that runs separately from the server.
// It is used for constant calculations, because it constantly runs tasks.
// There are two types of tasks:
// tasks - are executed asynchronously, but wait for all tasks in the iteration to be completed.
// notWaitTasks - are executed asynchronously, but do not wait for all tasks in the iteration to be completed.
// That is, task1 may still be executed for the first time, and task2 may already be executed 4 times.
// It is impossible to launch two instances of the same task, this means that task1 will only start again when it
// has completed its execution.
//
// IMPORTANT: it is recommended to set a delay of at least 1000 milliseconds (1 second).
type GlobalFlow struct {
	tasks             []Task
	noWaitTasks       sync.Map
	noWaitTasksIndex  int
	delayMilliseconds int
	working           sync.Map
}

func NewGlobalFlow(delayMilliseconds int) *GlobalFlow {
	return &GlobalFlow{
		noWaitTasksIndex:  0,
		delayMilliseconds: delayMilliseconds,
	}
}

// AddTask adding tasks that complete execution synchronously.
func (gf *GlobalFlow) AddTask(task Task) {
	gf.tasks = append(gf.tasks, task)
}

// AddNotWaitTask adds tasks that do not wait for synchronous completion.
func (gf *GlobalFlow) AddNotWaitTask(task Task) {
	gf.noWaitTasks.Store(gf.noWaitTasksIndex, task)
	gf.working.Store(gf.noWaitTasksIndex, false)
	gf.noWaitTasksIndex++
}

// Run starts the execution of two types of tasks in two separate goroutines.
func (gf *GlobalFlow) Run(manager interfaces.Manager) {
	if !typeopr.IsPointer(manager) {
		panic("The managerConfig must be passed by pointer.")
	}
	// Wait tasks.
	if len(gf.tasks) > 0 {
		var wg sync.WaitGroup
		go func() {
			for {
				for i := 0; i < len(gf.tasks); i++ {
					wg.Add(1)
					go func(i int) {
						defer wg.Done()
						gf.tasks[i](manager)
					}(i)
				}
				wg.Wait()
				time.Sleep(time.Duration(gf.delayMilliseconds) * time.Millisecond)
			}
		}()
	}
	// No wait tasks.
	if gf.noWaitTasksIndex > 0 {
		go func() {
			// Description of the algorithm:
			// 1. gf.working[i] is false before starting goroutines.
			// 2. gf.working[i] is true when the goroutine has started
			// 3. gf.working[i] is false when the goroutine has completed
			// 4. if gf.working[i] is true, the task is running and you should skip running this task.
			for {
				for i := 0; i < gf.noWaitTasksIndex; i++ {
					wrk, _ := gf.working.Load(i)
					if !wrk.(bool) {
						go func(i int) {
							defer func() {
								gf.working.Store(i, false)
							}()
							gf.working.Store(i, true)
							callTask, _ := gf.noWaitTasks.Load(i)
							callTask.(Task)(manager)
						}(i)
					}
				}
				time.Sleep(time.Duration(gf.delayMilliseconds) * time.Millisecond)
			}
		}()
	}
}
