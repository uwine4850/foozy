package database

import (
	"errors"
	"sync"

	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
)

type AsyncQueries struct {
	syncQ    interfaces.ISyncQueries
	wg       sync.WaitGroup
	asyncRes sync.Map
}

func NewAsyncQueries() *AsyncQueries {
	return &AsyncQueries{}
}

func (q *AsyncQueries) SetSyncQueries(queries interfaces.ISyncQueries) {
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
func AsyncResError(keys []string, db *Database) error {
	for i := 0; i < len(keys); i++ {
		res, ok := db.AsyncQ().LoadAsyncRes(keys[i])
		if !ok {
			return errors.New("key for loading the result of asynchronous query not found")
		}
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}
