package database

import (
	"fmt"
	"slices"
)

// ControlConnect manages connections to the database. You can only use one database instance per connection.
type ControlConnect struct {
	openConnections      []*Database
	openNamedConnections map[string]*Database
}

func (cc *ControlConnect) GetOpenUnnamedConnections() []*Database {
	return cc.openConnections
}

func (cc *ControlConnect) GetOpenNamedConnections() map[string]*Database {
	return cc.openNamedConnections
}

// OpenUnnamedConnection opens one unnamed connection to the database.
func (cc *ControlConnect) OpenUnnamedConnection(db *Database) error {
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

// CloseUnnamedConnectionByIndex closes unnamed connections by id.
func (cc *ControlConnect) CloseUnnamedConnectionByIndex(index int) error {
	if index >= 0 && index < len(cc.openConnections) {
		if err := cc.openConnections[index].Close(); err != nil {
			return err
		}
		cc.openConnections = append(cc.openConnections[:index], cc.openConnections[index+1:]...)
		return nil
	}
	return ErrConnectionNotExists{}
}

// CloseAllUnnamedConnection closes all unnamed connections.
func (cc *ControlConnect) CloseAllUnnamedConnection() error {
	for i := 0; i < len(cc.openConnections); i++ {
		err := cc.openConnections[i].Close()
		if err != nil {
			return err
		}
	}
	cc.openConnections = []*Database{}
	return nil
}

func (cc *ControlConnect) makeNamedConnectionMap() {
	if cc.openNamedConnections == nil {
		cc.openNamedConnections = make(map[string]*Database)
	}
}

// OpenNamedConnection opens a named connection.
func (cc *ControlConnect) OpenNamedConnection(name string, db *Database) error {
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

// CloseNamedConnection closes a named connection by name.
func (cc *ControlConnect) CloseNamedConnection(name string) error {
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

// CloseAllNamedConnection closes the entire named connection.
func (cc *ControlConnect) CloseAllNamedConnection() error {
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
