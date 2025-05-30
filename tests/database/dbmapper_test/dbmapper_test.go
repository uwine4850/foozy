package dbmappertest

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/uwine4850/foozy/pkg/database"
	qb "github.com/uwine4850/foozy/pkg/database/querybuld"
	"github.com/uwine4850/foozy/pkg/mapper"
	"github.com/uwine4850/foozy/pkg/typeopr"
	"github.com/uwine4850/foozy/tests/common/tconf"
)

var db = database.NewDatabase(tconf.DbArgs)

func createDbTest() {
	_, err := db.SyncQ().Query("INSERT INTO `dbtest` (`col1`, `col2`, `col3`) VALUES (?, ?, ?)",
		"test1", "2023-11-15", 111.22)
	if err != nil {
		panic(err)
	}
	_, err = db.SyncQ().Query("INSERT INTO `dbtest` (`col1`, `col2`, `col3`) VALUES (?, ?, ?)",
		"test2", "2023-11-20", 222.11)
	if err != nil {
		panic(err)
	}
}

func clearDbTest() {
	_, err := db.SyncQ().Query("TRUNCATE TABLE dbtest")
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	err := db.Connect()
	if err != nil {
		panic(err)
	}
	clearDbTest()
	createDbTest()
	exitCode := m.Run()
	err = db.Close()
	if err != nil {
		panic(err)
	}
	os.Exit(exitCode)
}

type DbTestMapper struct {
	Col1 string    `db:"col1"`
	Col2 time.Time `db:"col2"`
	Col3 float64   `db:"col3"`
	Col4 int       `db:"col4" empty:"0"`
}

func TestDbMapperUseStruct(t *testing.T) {
	clearDbTest()
	createDbTest()
	q := qb.NewSyncQB(db.SyncQ()).SelectFrom("*", "dbtest")
	q.Merge()
	res, err := q.Query()
	if err != nil {
		t.Error(err)
	}
	dbTestMapper := make([]DbTestMapper, len(res))
	if err := mapper.FillStructSliceFromDb(&dbTestMapper, &res); err != nil {
		t.Error(err)
	}
	if len(dbTestMapper) == 0 {
		t.Error("DbMapper.Output must not be empty")
	}
	tm, err := time.Parse(time.DateOnly, "2023-11-15")
	if err != nil {
		t.Error(err)
	}
	map1 := DbTestMapper{Col1: "test1", Col2: tm, Col3: 111.22, Col4: 0}
	if dbTestMapper[0] != map1 {
		t.Error("DbMapper.Output value does not match expected")
	}
	tm2, err := time.Parse(time.DateOnly, "2023-11-20")
	if err != nil {
		t.Error(err)
	}
	map2 := DbTestMapper{Col1: "test2", Col2: tm2, Col3: 222.11, Col4: 0}
	if dbTestMapper[1] != map2 {
		t.Error("DbMapper.Output value does not match expected")
	}
}

// func TestDbMapperUseMap(t *testing.T) {
// 	clearDbTest()
// 	createDbTest()
// 	q := qb.NewSyncQB(db.SyncQ()).SelectFrom("*", "dbtest")
// 	q.Merge()
// 	res, err := q.Query()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	var dbTestMapper = []map[string]string{}
// 	mapper := dbmapper.NewMapper(res, typeopr.Ptr{}.New(&dbTestMapper))
// 	if err := mapper.Fill(); err != nil {
// 		t.Error(err)
// 	}
// 	if len(dbTestMapper) == 0 {
// 		t.Error("DbMapper.Output must not be empty")
// 	}
// 	map1 := map[string]string{"col1": "test1", "col2": "2023-11-15", "col3": "111.22", "col4": "0", "id": ""}
// 	if !fmap.Compare(&map1, &dbTestMapper[0], []string{"id"}) {
// 		t.Error("DbMapper.Output value does not match expected")
// 	}
// 	map2 := map[string]string{"col1": "test2", "col2": "2023-11-20", "col3": "222.11", "col4": "0", "id": ""}
// 	if !fmap.Compare(&map2, &dbTestMapper[1], []string{"id"}) {
// 		t.Error("DbMapper.Output value does not match expected")
// 	}
// }

type Fill struct {
	Col1 string    `db:"col1"`
	Col2 time.Time `db:"col2"`
	Col3 float64   `db:"col3"`
}

func TestFillStructFromDb(t *testing.T) {
	tm, err := time.Parse(time.DateOnly, "2023-11-15")
	if err != nil {
		t.Error(err)
	}
	expected := Fill{
		Col1: "test1",
		Col2: tm,
		Col3: 111.22,
	}
	res, err := db.SyncQ().Query("SELECT * FROM dbtest")
	if err != nil {
		t.Error(err)
	}
	var f Fill
	err = mapper.FillStructFromDb(typeopr.Ptr{}.New(&f), &res[0])
	if err != nil {
		t.Error(err)
	}
	if f != expected {
		t.Errorf("The data in the structure is not as expected.")
	}
}

func TestFillReflectValueFromDb(t *testing.T) {
	tm, err := time.Parse(time.DateOnly, "2023-11-15")
	if err != nil {
		t.Error(err)
	}
	expected := Fill{
		Col1: "test1",
		Col2: tm,
		Col3: 111.22,
	}
	res, err := db.SyncQ().Query("SELECT * FROM dbtest")
	if err != nil {
		t.Error(err)
	}
	f := Fill{}
	fValue := reflect.ValueOf(&f).Elem()
	err = mapper.FillStructFromDb(typeopr.Ptr{}.New(&fValue), &res[0])
	if err != nil {
		t.Error(err)
	}
	if f != expected {
		t.Errorf("The data in the structure is not as expected.")
	}
}

// func TestFillMapFromDb(t *testing.T) {
// 	expected := map[string]string{"col1": "test1", "col2": "2023-11-15", "col3": "111.22"}
// 	res, err := db.SyncQ().Query("SELECT * FROM dbtest")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	m := map[string]string{}
// 	if err := dbmapper.FillMapFromDb(res[0], &m); err != nil {
// 		t.Error(err)
// 	}
// 	if !fmap.Compare(&expected, &m, []string{"id", "col4"}) {
// 		t.Error("the completed map does not match the expected one")
// 	}
// }
