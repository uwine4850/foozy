## builtin globalflow

#### KeyUpdater
Function for use in __TODO:add link__ [GlobalFlow.AddNotWaitTask]().<br>
Updates the hashKey and blockKey keys after a selected period of time.
```golang
func KeyUpdater(delaySec int) globalflow.Task {
	return func(manager interfaces.Manager) {
		time.Sleep(time.Duration(delaySec) * time.Second)
		manager.Key().Get32BytesKey().GenerateBytesKeys(32)
	}
}
```