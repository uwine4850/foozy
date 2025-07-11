package database

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
)

type SyncQueries struct {
	qe interfaces.QueryExec
}

func NewSyncQueries() *SyncQueries {
	return &SyncQueries{}
}

func (q *SyncQueries) New() (interface{}, error) {
	return &SyncQueries{
		qe: q.qe,
	}, nil
}

// Query wrapper for the IDbQuery.Query method.
func (q *SyncQueries) Query(query string, args ...any) ([]map[string]interface{}, error) {
	return q.qe.Query(query, args...)
}

// Exec wrapper for the IDbQuery.Exec method.
func (q *SyncQueries) Exec(query string, args ...any) (map[string]interface{}, error) {
	return q.qe.Exec(query, args...)
}

func (q *SyncQueries) SetDB(qe interfaces.QueryExec) {
	q.qe = qe
}
