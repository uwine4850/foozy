package globalflow

import (
	"sync"
	"time"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/utils"
)

type Task func(manager interfaces.IManager)

type GlobalFlow struct {
	tasks []Task
	delay time.Duration
}

func NewGlobalFlow(delay time.Duration) *GlobalFlow {
	return &GlobalFlow{delay: delay}
}

func (gf *GlobalFlow) AddTask(task Task) {
	gf.tasks = append(gf.tasks, task)
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
					gf.tasks[i](manager)
					wg.Done()
				}(i)
			}
			wg.Wait()
			time.Sleep(gf.delay)
		}
	}()
}
