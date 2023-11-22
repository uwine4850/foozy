package database

import (
	"github.com/uwine4850/foozy/pkg/database/dbutils"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"sync"
)

type AsyncQueries struct {
	syncQ    interfaces.ISyncQueries
	wg       sync.WaitGroup
	asyncRes sync.Map
	interfaces.IQueryBuild
}

func NewAsyncQueries(qb interfaces.IQueryBuild) *AsyncQueries {
	return &AsyncQueries{IQueryBuild: qb}
}

func (q *AsyncQueries) QB(key string) interfaces.IQueryBuild {
	q.IQueryBuild.SetAsyncQ(q)
	q.IQueryBuild.SetKeyForAsyncQ(key)
	return q.IQueryBuild
}

func (q *AsyncQueries) SetSyncQueries(queries interfaces.ISyncQueries) {
	q.syncQ = queries
}

func (q *AsyncQueries) AsyncQuery(key string, query string, args ...any) {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		_res, err := q.syncQ.Query(query, args...)
		q.setAsyncRes(key, _res, err)
	}()
}

func (q *AsyncQueries) AsyncSelect(key string, rows []string, tableName string, where dbutils.WHOutput, limit int) {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		_res, err := q.syncQ.Select(rows, tableName, where, limit)
		q.setAsyncRes(key, _res, err)
	}()
}

func (q *AsyncQueries) AsyncInsert(key string, tableName string, params map[string]interface{}) {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		_res, err := q.syncQ.Insert(tableName, params)
		q.setAsyncRes(key, _res, err)
	}()
}

func (q *AsyncQueries) AsyncUpdate(key string, tableName string, params []dbutils.DbEquals, where dbutils.WHOutput) {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		_res, err := q.syncQ.Update(tableName, params, where)
		q.setAsyncRes(key, _res, err)
	}()
}

func (q *AsyncQueries) AsyncDelete(key string, tableName string, where dbutils.WHOutput) {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		_res, err := q.syncQ.Delete(tableName, where)
		q.setAsyncRes(key, _res, err)
	}()
}

func (q *AsyncQueries) AsyncCount(key string, rows []string, tableName string, where dbutils.WHOutput, limit int) {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		_res, err := q.syncQ.Count(rows, tableName, where, limit)
		q.setAsyncRes(key, _res, err)
	}()
}

// setAsyncRes sets the result of the key command execution.
func (q *AsyncQueries) setAsyncRes(key string, _res []map[string]interface{}, err error) {
	queryData := dbutils.AsyncQueryData{}
	if err != nil {
		queryData.Error = err
	}
	queryData.Res = _res
	q.asyncRes.Store(key, queryData)
}

// LoadAsyncRes retrieves command execution data by key.
func (q *AsyncQueries) LoadAsyncRes(key string) (*dbutils.AsyncQueryData, bool) {
	value, ok := q.asyncRes.Load(key)
	v := value.(dbutils.AsyncQueryData)
	return &v, ok
}

// Wait waits for the execution of all asynchronous methods that are started before executing this method.
// IMPOrTANT: this method must be run before LoadAsyncRes.
// Several Wait methods can be called if necessary.
func (q *AsyncQueries) Wait() {
	q.wg.Wait()
}
