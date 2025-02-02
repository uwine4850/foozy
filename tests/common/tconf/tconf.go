package tconf

import "github.com/uwine4850/foozy/pkg/database"

const (
	PortAuth        = ":5000"
	PortDebug       = ":5001"
	PortForm        = ":5002"
	PortFormMapping = ":5003"
	PortMddl        = ":5004"
	PortObject      = ":5005"
	PortRestMapper  = ":5006"
	PortRouter      = ":5007"
	PortCookies     = ":5008"
)

var DbArgs = database.DbArgs{
	Username: "root", Password: "1111", Host: "localhost", Port: "3408", DatabaseName: "foozy_test",
}
