## database

### AsyncQueries
asynchronous database queries.

The `ISyncQueries` is used to implement the queries, only it is wrapped in this object.
For each query, you must specify a key by which it can be identified.

__IMPORTANT_1__: it is necessary to call [Wait] method to correctly wait for queries execution.
__IMPORTANT_2__: for each new asynchronous request you must create a separate instance of this object.
This is done to protect the user from data leakage, because the object saves user request data and should not be shared.
```golang
type AsyncQueries struct {
	syncQ    interfaces.SyncQ
	wg       sync.WaitGroup
	asyncRes sync.Map
}
```

#### AsyncQueries.LoadAsyncRes
Retrieves command execution data by key.
```golang
func (q *AsyncQueries) LoadAsyncRes(key string) (*dbutils.AsyncQueryData, bool) {
	value, ok := q.asyncRes.Load(key)
	if ok {
		v := value.(*dbutils.AsyncQueryData)
		return v, ok
	}
	return nil, ok
}
```

#### AsyncQueries.Wait
Waits for the execution of all asynchronous methods that are started before executing this method.

__IMPORTANT__: this method must be run before LoadAsyncRes.

Several Wait methods can be called if necessary.
```golang
func (q *AsyncQueries) Wait() {
	q.wg.Wait()
}
```

#### AsyncQueries.Clear
Clears the query results data.
```golang
func (q *AsyncQueries) Clear() {
	q.asyncRes = sync.Map{}
}
```

#### AsyncQueries.AsyncResError
Loads the result of several asynchronous key queries and checks for errors.
```golang
func AsyncResError(keys []string, asyncQ interfaces.AsyncQ) error {
	for i := 0; i < len(keys); i++ {
		res, ok := asyncQ.LoadAsyncRes(keys[i])
		if !ok {
			return errors.New("key for loading the result of asynchronous query not found")
		}
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}
```