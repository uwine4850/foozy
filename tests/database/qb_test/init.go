package qbtest

import (
	"errors"

	"github.com/uwine4850/foozy/pkg/database"
	qb "github.com/uwine4850/foozy/pkg/database/querybuld"
	"github.com/uwine4850/foozy/tests/common/tconf"
)

var db = database.NewDatabase(tconf.DbArgs)

// var syncQB *qb.QB = qb.NewSyncQB(db.SyncQ())

var testErrorText = errors.New("the query text does not match the expected text")
var testErrorArgs = errors.New("arguments do not match expectations")

func syncQB() *qb.QB {
	return qb.NewSyncQB(db.SyncQ())
}
