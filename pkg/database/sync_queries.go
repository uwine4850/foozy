package database

import (
	"github.com/uwine4850/foozy/pkg/interfaces"
)

type SyncQueries struct {
	db interfaces.IDbQuery
}

func NewSyncQueries() *SyncQueries {
	return &SyncQueries{}
}

// Query wrapper for the IDbQuery.Query method.
func (q *SyncQueries) Query(query string, args ...any) ([]map[string]interface{}, error) {
	return q.db.Query(query, args...)
}

// Exec wrapper for the IDbQuery.Exec method.
func (q *SyncQueries) Exec(query string, args ...any) (map[string]interface{}, error) {
	return q.db.Exec(query, args...)
}

func (q *SyncQueries) SetDB(db interfaces.IDbQuery) {
	q.db = db
}
