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
	qb       interfaces.IQueryBuild
}

func NewAsyncQueries(qb interfaces.IQueryBuild) *AsyncQueries {
	return &AsyncQueries{qb: qb}
}

func (q *AsyncQueries) QB(key string) interfaces.IUserQueryBuild {
	q.qb.SetAsyncQ(q)
	q.qb.SetKeyForAsyncQ(key)
	return q.qb
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

func (q *AsyncQueries) Select(key string, rows []string, tableName string, where dbutils.WHOutput, limit int) {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		res, err := q.syncQ.Select(rows, tableName, where, limit)
		q.storeAsyncRes(key, &dbutils.AsyncQueryData{Res: res, Error: err})
	}()
}

func (q *AsyncQueries) Insert(key string, tableName string, params map[string]any) {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		res, err := q.syncQ.Insert(tableName, params)
		q.storeAsyncRes(key, &dbutils.AsyncQueryData{SingleRes: res, Error: err})
	}()
}

func (q *AsyncQueries) Update(key string, tableName string, params map[string]interface{}, where dbutils.WHOutput) {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		res, err := q.syncQ.Update(tableName, params, where)
		q.storeAsyncRes(key, &dbutils.AsyncQueryData{SingleRes: res, Error: err})
	}()
}

func (q *AsyncQueries) Delete(key string, tableName string, where dbutils.WHOutput) {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		res, err := q.syncQ.Delete(tableName, where)
		q.storeAsyncRes(key, &dbutils.AsyncQueryData{SingleRes: res, Error: err})
	}()
}

func (q *AsyncQueries) Count(key string, rows []string, tableName string, where dbutils.WHOutput, limit int) {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		res, err := q.syncQ.Count(rows, tableName, where, limit)
		q.storeAsyncRes(key, &dbutils.AsyncQueryData{Res: res, Error: err})
	}()
}

func (q *AsyncQueries) Increment(key string, fieldName string, tableName string, where dbutils.WHOutput) {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		res, err := q.syncQ.Increment(fieldName, tableName, where)
		q.storeAsyncRes(key, &dbutils.AsyncQueryData{SingleRes: res, Error: err})
	}()
}

func (q *AsyncQueries) Exists(key string, tableName string, where dbutils.WHOutput) {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		res, err := q.syncQ.Exists(tableName, where)
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
		q.asyncRes.Delete(key)
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
