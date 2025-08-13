## package globalflow

### GlobalFlow
Creates a flow that runs separately from the server.<br>
It is used for constant calculations, because it constantly runs tasks.

There are two types of tasks:

* tasks - are executed asynchronously, but wait for all tasks in the iteration to be completed.
* notWaitTasks - are executed asynchronously, but do not wait for all tasks in the iteration to be completed.

That is, `task1` may still be executed for the first time, and `task2` may already be executed 4 times.
It is impossible to launch two instances of the same task, this means that `task1` will only start again when it 
has completed its execution.

__IMPORTANT:__ it is recommended to set a delay of at least 1000 milliseconds (1 second).

#### GlobalFlow.AddTask
Adding tasks that complete execution synchronously.
```golang
func (gf *GlobalFlow) AddTask(task Task) {
	gf.tasks = append(gf.tasks, task)
}
```

#### GlobalFlow.AddNotWaitTask
Adds tasks that do not wait for synchronous completion.
```golang
func (gf *GlobalFlow) AddNotWaitTask(task Task) {
	gf.noWaitTasks.Store(gf.noWaitTasksIndex, task)
	gf.working.Store(gf.noWaitTasksIndex, false)
	gf.noWaitTasksIndex++
}
```

#### GlobalFlow.Run
Starts the execution of two types of tasks in two separate goroutines.
```golang
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

```