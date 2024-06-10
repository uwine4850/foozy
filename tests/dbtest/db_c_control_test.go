package dbtest

import (
	"errors"
	"strconv"
	"testing"

	"github.com/uwine4850/foozy/pkg/database"
)

var dbcc = database.NewDatabase("root", "1111", "localhost", "3408", "foozy_test")

func TestOpenCloseAllConnections(t *testing.T) {
	cc := database.ConnectControl{}
	if err := cc.OpenConnection(dbcc); err != nil {
		t.Error(err)
	}
	if err := cc.OpenConnection(dbcc); err != nil {
		t.Error(err)
	}
	if err := cc.OpenConnection(dbcc); err != nil {
		t.Error(err)
	}
	if len(cc.GetOpenConnections()) != 3 {
		t.Errorf("Error opening connections. There should be 3, there are %s.", strconv.Itoa(len(cc.GetOpenConnections())))
	}
	err := cc.CloseAllConnection()
	if err != nil {
		t.Error(err)
	}
	if len(cc.GetOpenConnections()) != 0 {
		t.Errorf("Error closing connections, %s connections left.", strconv.Itoa(len(cc.GetOpenConnections())))
	}
}

func TestOpenCloseNamedConnections(t *testing.T) {
	cc := database.ConnectControl{}
	if err := cc.OpenNamedConnection("1", dbcc); err != nil {
		t.Error(err)
	}
	if err := cc.OpenNamedConnection("2", dbcc); err != nil {
		t.Error(err)
	}
	if err := cc.OpenNamedConnection("3", dbcc); err != nil {
		t.Error(err)
	}
	if len(cc.GetOpenNamedConnections()) != 3 {
		t.Errorf("Error opening named connections. There should be 3, there are %s.", strconv.Itoa(len(cc.GetOpenNamedConnections())))
	}
	if err := cc.CloseNamedConnection("1"); err != nil {
		t.Error(err)
	}
	if err := cc.CloseNamedConnection("2"); err != nil {
		t.Error(err)
	}
	if err := cc.CloseNamedConnection("3"); err != nil {
		t.Error(err)
	}
	if len(cc.GetOpenConnections()) != 0 {
		t.Errorf("Error closing named connections, %s connections left.", strconv.Itoa(len(cc.GetOpenNamedConnections())))
	}
}

func TestNamedConnectionOpenCloseError(t *testing.T) {
	cc := database.ConnectControl{}
	if err := cc.OpenNamedConnection("1", dbcc); err != nil {
		t.Error(err)
	}
	okErr := cc.OpenNamedConnection("1", dbcc)
	openError := database.ErrNamedConnectionAlreadyExists{ConnectionName: "1"}
	if !errors.Is(okErr, openError) {
		t.Errorf("The name was repeated, but the error was not displayed.")
	}
	if err := cc.CloseNamedConnection("1"); err != nil {
		t.Error(err)
	}
	okErr = cc.CloseNamedConnection("1")
	closeErr := database.ErrNamedConnectionNotExists{ConnectionName: "1"}
	if !errors.Is(okErr, closeErr) {
		t.Errorf("The name does not exist, but no error was displayed.")
	}
}

func TestNamedConnectionCloseAll(t *testing.T) {
	cc := database.ConnectControl{}
	if err := cc.OpenNamedConnection("1", dbcc); err != nil {
		t.Error(err)
	}
	if err := cc.OpenNamedConnection("2", dbcc); err != nil {
		t.Error(err)
	}
	if err := cc.OpenNamedConnection("3", dbcc); err != nil {
		t.Error(err)
	}
	err := cc.CloseAllConnection()
	if err != nil {
		t.Error(err)
	}
	if len(cc.GetOpenConnections()) != 0 {
		t.Errorf("Error closing named connections, %s connections left.", strconv.Itoa(len(cc.GetOpenNamedConnections())))
	}
}
