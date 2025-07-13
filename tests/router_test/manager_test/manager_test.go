package manager_test

import (
	"errors"
	"testing"

	"github.com/uwine4850/foozy/pkg/database"
	"github.com/uwine4850/foozy/pkg/interfaces"
	"github.com/uwine4850/foozy/pkg/router/manager"
)

var newManager interfaces.Manager

func TestMain(m *testing.M) {
	newOTD := manager.NewOneTimeData()
	newDBPool := database.NewDatabasePool()
	newManager = manager.NewManager(newOTD, nil, newDBPool)
	m.Run()
}

func TestOneTimeData(t *testing.T) {
	newManager.OneTimeData().SetUserContext("key", true)
	key, isKey := newManager.OneTimeData().GetUserContext("key")
	if !isKey {
		t.Error("value not found")
	}
	_, ok := key.(bool)
	if !ok {
		t.Error("value does not match the expectation")
	}
	newManager.OneTimeData().DelUserContext("key")
	_, isKey = newManager.OneTimeData().GetUserContext("key")
	if isKey {
		t.Error("deleted value found")
	}
	newManager.OneTimeData().SetSlugParams(map[string]string{"id": "1"})
	_, isSlug := newManager.OneTimeData().GetSlugParams("id")
	if !isSlug {
		t.Error("slug parameter not found")
	}
}

type fakeDatabase struct{}

func (d *fakeDatabase) SyncQ() interfaces.SyncQ {
	return nil
}
func (d *fakeDatabase) NewAsyncQ() (interfaces.AsyncQ, error) {
	return nil, nil
}
func (d *fakeDatabase) NewTransaction() (interfaces.DatabaseTransaction, error) {
	return nil, nil
}

func TestConnectionPool(t *testing.T) {
	newManager.Database().AddConnection("pool", &fakeDatabase{})
	_, err := newManager.Database().ConnectionPool("pool")
	if !errors.Is(err, &database.ErrDatabasePoolNotLocked{}) {
		t.Error("expected error not received")
	}
	newManager.Database().Lock()
	_, err = newManager.Database().ConnectionPool("pool")
	if err != nil {
		t.Error(err)
	}
	err = newManager.Database().AddConnection("pool", &fakeDatabase{})
	if !errors.Is(err, &database.ErrDatabasePoolIsLocked{}) {
		t.Error("expected error not received")
	}
}
