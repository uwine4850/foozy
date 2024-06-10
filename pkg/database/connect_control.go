package database

import (
	"fmt"
	"slices"
)

type ConnectControl struct {
	openConnections      []*Database
	openNamedConnections map[string]*Database
}

func (cc *ConnectControl) GetOpenConnections() []*Database {
	return cc.openConnections
}

func (cc *ConnectControl) GetOpenNamedConnections() map[string]*Database {
	return cc.openNamedConnections
}

func (cc *ConnectControl) OpenConnection(db *Database) error {
	if slices.Contains(cc.openConnections, db) {
		return ErrConnectionAlreadyExists{}
	}
	if err := db.Connect(); err != nil {
		return err
	} else {
		cc.openConnections = append(cc.openConnections, db)
		return nil
	}
}

func (cc *ConnectControl) CloseConnectionByIndex(index int) error {
	if index >= 0 && index < len(cc.openConnections) {
		if err := cc.openConnections[index].Close(); err != nil {
			return err
		}
		cc.openConnections = append(cc.openConnections[:index], cc.openConnections[index+1:]...)
		return nil
	}
	return ErrConnectionNotExists{}
}

func (cc *ConnectControl) CloseAllConnection() error {
	for i := 0; i < len(cc.openConnections); i++ {
		err := cc.openConnections[i].Close()
		if err != nil {
			return err
		}
	}
	cc.openConnections = []*Database{}
	return nil
}

func (cc *ConnectControl) makeNamedConnectionMap() {
	if cc.openNamedConnections == nil {
		cc.openNamedConnections = make(map[string]*Database)
	}
}

func (cc *ConnectControl) OpenNamedConnection(name string, db *Database) error {
	_, ok := cc.openNamedConnections[name]
	if ok {
		return ErrNamedConnectionAlreadyExists{name}
	}
	if err := db.Ping(); err == nil {
		return ErrConnectionAlreadyOpen{}
	}
	if err := db.Connect(); err != nil {
		return err
	} else {
		cc.makeNamedConnectionMap()
		cc.openNamedConnections[name] = db
		return nil
	}
}

func (cc *ConnectControl) CloseNamedConnection(name string) error {
	cc.makeNamedConnectionMap()
	_, ok := cc.openNamedConnections[name]
	if ok {
		conn := cc.openNamedConnections[name]
		if err := conn.Close(); err != nil {
			return err
		}
		delete(cc.openNamedConnections, name)
		return nil
	} else {
		return ErrNamedConnectionNotExists{name}
	}
}

func (cc *ConnectControl) CloseAllNamedConnection() error {
	cc.makeNamedConnectionMap()
	for _, conn := range cc.openNamedConnections {
		if err := conn.Close(); err != nil {
			return err
		}
	}
	cc.openNamedConnections = make(map[string]*Database)
	return nil
}

type ErrConnectionAlreadyOpen struct {
}

func (e ErrConnectionAlreadyOpen) Error() string {
	return "Connection already open."
}

type ErrConnectionAlreadyExists struct {
	ConnectionName string
}

func (e ErrConnectionAlreadyExists) Error() string {
	return "Connection already exists."
}

type ErrConnectionNotExists struct {
	ConnectionName string
}

func (e ErrConnectionNotExists) Error() string {
	return "Connection not exists."
}

type ErrNamedConnectionAlreadyExists struct {
	ConnectionName string
}

func (e ErrNamedConnectionAlreadyExists) Error() string {
	return fmt.Sprintf("A named connection \"%s\" already exists.", e.ConnectionName)
}

type ErrNamedConnectionNotExists struct {
	ConnectionName string
}

func (e ErrNamedConnectionNotExists) Error() string {
	return fmt.Sprintf("A named connection \"%s\" not exists.", e.ConnectionName)
}
