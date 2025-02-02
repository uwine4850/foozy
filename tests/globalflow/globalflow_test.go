package globalflowtest

import (
	"os"
	"testing"
	"time"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router/manager"
	"github.com/uwine4850/foozy/pkg/server/globalflow"
)

var mng = manager.NewManager(nil)

func TestMain(m *testing.M) {
	mng.SetOneTimeData(manager.NewManagerData())
	gf := globalflow.NewGlobalFlow(1)
	gf.AddTask(func(manager interfaces.IManager) {
		mng.OneTimeData().SetUserContext("TASK_1", "TASK 1")
	})
	gf.AddTask(func(manager interfaces.IManager) {
		mng.OneTimeData().SetUserContext("TASK_2", "TASK 2")
	})
	gf.AddNotWaitTask(func(manager interfaces.IManager) {
		mng.OneTimeData().SetUserContext("NOT_WAIT_TASK_1", "NOT WAIT TASK 1")
	})
	gf.AddNotWaitTask(func(manager interfaces.IManager) {
		mng.OneTimeData().SetUserContext("NOT_WAIT_TASK_2", "NOT WAIT TASK 2")
	})
	gf.Run(mng)
	time.Sleep(1 * time.Second)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestTask(t *testing.T) {
	time.Sleep(2 * time.Second)
	_, okT1 := mng.OneTimeData().GetUserContext("TASK_1")
	if !okT1 {
		t.Errorf("Task_1 failed")
	}
	_, okT2 := mng.OneTimeData().GetUserContext("TASK_2")
	if !okT2 {
		t.Errorf("Task_2 failed")
	}
}

func TestNotWaitTask(t *testing.T) {
	time.Sleep(2 * time.Second)
	_, okT1 := mng.OneTimeData().GetUserContext("NOT_WAIT_TASK_1")
	if !okT1 {
		t.Errorf("NOT_WAIT_TASK_1 failed")
	}
	_, okT2 := mng.OneTimeData().GetUserContext("NOT_WAIT_TASK_2")
	if !okT2 {
		t.Errorf("NOT_WAIT_TASK_2 failed")
	}
}
