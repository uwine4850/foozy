package globalflow

import (
	"sync"
	"time"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/utils"
)

type Task func(manager interfaces.IManager)

type GlobalFlow struct {
	tasks             []Task
	notWaitTasks      map[int]Task
	notWaitTasksIndex int
	delay             time.Duration
	working           map[int]bool
	mutex             sync.Mutex
}

func NewGlobalFlow(delay time.Duration) *GlobalFlow {
	return &GlobalFlow{
		notWaitTasks:      make(map[int]Task),
		notWaitTasksIndex: 0,
		delay:             delay,
		working:           make(map[int]bool),
	}
}

func (gf *GlobalFlow) AddTask(task Task) {
	gf.tasks = append(gf.tasks, task)
}

func (gf *GlobalFlow) AddNotWaitTask(task Task) {
	gf.notWaitTasks[gf.notWaitTasksIndex] = task
	gf.working[gf.notWaitTasksIndex] = false
	gf.notWaitTasksIndex++
}

func (gf *GlobalFlow) Run(manager interfaces.IManager) {
	if !utils.IsPointer(manager) {
		panic("The manager must be passed by pointer.")
	}
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
			time.Sleep(gf.delay)
		}
	}()
	go func() {
		// Description of the algorithm:
		// 1. gf.working[i] is false before starting goroutines.
		// 2. gf.working[i] is true when the goroutine has started
		// 3. gf.working[i] is false when the goroutine has completed
		// 4. if gf.working[i] is true, the task is running and you should skip running this task.
		for {
			for i := 0; i < gf.notWaitTasksIndex; i++ {
				if !gf.working[i] {
					go func(i int) {
						defer func() {
							gf.mutex.Lock()
							gf.working[i] = false // Task done.
							gf.mutex.Unlock()
						}()
						gf.mutex.Lock()
						gf.working[i] = true // Task is running.
						gf.mutex.Unlock()
						gf.notWaitTasks[i](manager)
					}(i)
				}
			}
			time.Sleep(gf.delay)
		}
	}()
}
