package database

import (
	"fmt"
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
	if err := db.Connect(); err != nil {
		return err
	} else {
		cc.openConnections = append(cc.openConnections, db)
		return nil
	}
}

func (cc *ConnectControl) makeNamedConnectionMap() {
	if cc.openNamedConnections == nil {
		cc.openNamedConnections = make(map[string]*Database)
	}
}

func (cc *ConnectControl) OpenNamedConnection(name string, db *Database) error {
	if err := db.Connect(); err != nil {
		return err
	} else {
		cc.makeNamedConnectionMap()
		_, ok := cc.openNamedConnections[name]
		if ok {
			return ErrNamedConnectionAlreadyExists{name}
		} else {
			cc.openNamedConnections[name] = db
		}
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
	for name, conn := range cc.openNamedConnections {
		if err := conn.Close(); err != nil {
			return err
		}
		delete(cc.openNamedConnections, name)
	}
	return nil
}

func (cc *ConnectControl) CloseAllConnection() error {
	for i := 0; i < len(cc.openConnections); i++ {
		err := cc.openConnections[i].Close()
		if err != nil {
			return err
		} else {
			if i == len(cc.openConnections)-1 {
				cc.openConnections = []*Database{}
			} else {
				cc.openConnections = append(cc.openConnections[:i], cc.openConnections[i+1:]...)
			}
		}
	}
	return nil
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
