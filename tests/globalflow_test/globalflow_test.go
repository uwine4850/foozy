package globalflow_test

import (
	"sync"
	"testing"

	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/server/globalflow"
)

// sync.WaitGroup is only used in tests.
// It is necessary to prevent the test from terminating before the scheduled time.
func TestGlobalflow(t *testing.T) {
	noWaitTask := false
	waitTask := false
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Add(1)

	mngr := manager.NewManager(manager.NewOneTimeData(), nil, database.NewDatabasePool())
	gf := globalflow.NewGlobalFlow(1000)

	gf.AddTask(func(manager interfaces.Manager) {
		defer wg.Done()
		if !noWaitTask {
			t.Error("no wait task has not yet been completed")
		}
		waitTask = true
	})
	gf.AddNotWaitTask(func(manager interfaces.Manager) {
		defer wg.Done()
		noWaitTask = true
	})

	gf.Run(mngr)

	wg.Wait()
	if !waitTask {
		t.Error("wait task not yet completed")
	}
}
