package database

import (
	"errors"
	"sync"

	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
)

// AsyncQueries asynchronous database queries.
// The [ISyncQueries] is used to implement the queries, only it is wrapped in this object.
// For each query, you must specify a key by which it can be identified.
// IMPORTANT_1: it is necessary to call [Wait] method to correctly wait for queries execution.
// IMPORTANT_2: for each new asynchronous request you must create a separate instance of this object.
// This is done to protect the user from data leakage, because the object saves user request data
// and should not be shared.
type AsyncQueries struct {
	syncQ    interfaces.SyncQ
	wg       sync.WaitGroup
	asyncRes sync.Map
}

func NewAsyncQueries(syncQ interfaces.SyncQ) *AsyncQueries {
	return &AsyncQueries{
		syncQ: syncQ,
	}
}

func (q *AsyncQueries) New() (interface{}, error) {
	return &AsyncQueries{
		syncQ: q.syncQ,
	}, nil
}

func (q *AsyncQueries) SetSyncQueries(queries interfaces.SyncQ) {
	q.syncQ = queries
}

func (q *AsyncQueries) Query(key string, query string, args ...any) {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		res, err := q.syncQ.Query(query, args...)
		q.storeAsyncRes(key, &dbutils.AsyncQueryData{Res: res, Error: err})
	}()
}

func (q *AsyncQueries) Exec(key string, query string, args ...any) {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		res, err := q.syncQ.Exec(query, args...)
		q.storeAsyncRes(key, &dbutils.AsyncQueryData{SingleRes: res, Error: err})
	}()
}

// storeAsyncRes sets the result of the key command execution.
func (q *AsyncQueries) storeAsyncRes(key string, asyncQueryData *dbutils.AsyncQueryData) {
	q.asyncRes.Store(key, asyncQueryData)
}

// LoadAsyncRes retrieves command execution data by key.
func (q *AsyncQueries) LoadAsyncRes(key string) (*dbutils.AsyncQueryData, bool) {
	value, ok := q.asyncRes.Load(key)
	if ok {
		v := value.(*dbutils.AsyncQueryData)
		return v, ok
	}
	return nil, ok
}

// Wait waits for the execution of all asynchronous methods that are started before executing this method.
// IMPORTANT: this method must be run before LoadAsyncRes.
// Several Wait methods can be called if necessary.
func (q *AsyncQueries) Wait() {
	q.wg.Wait()
}

// Clear clears the query results data.
func (q *AsyncQueries) Clear() {
	q.asyncRes = sync.Map{}
}

// AsyncResError loads the result of several asynchronous key queries and checks for errors.
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
