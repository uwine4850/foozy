package bglobalflow

import (
	"time"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/server/globalflow"
)

// KeyUpdater function for use in GlobalFlow.AddNot WaitTask.
// Updates the hashKey and blockKey keys after a selected period of time.
func KeyUpdater(delaySec int) globalflow.Task {
	return func(manager interfaces.IManager) {
		time.Sleep(time.Duration(delaySec) * time.Second)
		manager.Key().Get32BytesKey().GenerateBytesKeys(32)
	}
}
