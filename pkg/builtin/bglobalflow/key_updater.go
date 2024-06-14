package bglobalflow

import (
	"time"

	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/server/globalflow"
)

func KeyUpdater(delaySec int) globalflow.Task {
	return func(manager interfaces.IManager) {
		time.Sleep(time.Duration(delaySec) * time.Second)
		manager.Get32BytesKey().GenerateBytesKeys(32)
	}
}
