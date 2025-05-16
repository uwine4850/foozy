package dbtest

import (
	"fmt"
	"testing"
)

func TestAsyncQuery(t *testing.T) {
	asyncQ := db.NewAsyncQ()
	asyncQ.Query("s", "SELECT `col1`, `col2`, `col3` FROM `db_async_test` LIMIT 1")
	asyncQ.Wait()
	res, _ := asyncQ.LoadAsyncRes("s")
	if res.Error != nil {
		t.Error(res.Error)
	}
	s := fmt.Sprintf("%v", res.Res)
	if s != "[map[col1:[116 101 115 116 49] col2:[50 48 50 51 45 49 49 45 49 53] col3:111.22]]" {
		t.Errorf("The row values in the database of the expected row do not match.")
	}
}
