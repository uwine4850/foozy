package qb

import "fmt"

// SelectExists checks if there is a value in the table.
// It is important to use a condition for correct operation.
func SelectExists(qb *QB, tableName string, whereValues ...any) (bool, error) {
	exSQ := SQ(false, NewNoDbQB().SelectFrom("1", tableName).Where(whereValues...))
	baseSQ := SQ(false, NewNoDbQB().Select(Exists(exSQ)).As("is_exists"))
	qb.Func(baseSQ).Merge()
	res, err := qb.Query()
	if err != nil {
		return false, err
	}
	return res[0]["is_exists"].(int64) != 0, nil
}

// Increment increases the numeric value of the table by one.
func Increment(qb *QB, tableName string, field string, whereValues ...any) (bool, error) {
	customQ := fmt.Sprintf("UPDATE %s SET %s = %s + 1", tableName, field, field)
	baseSQ := SQ(false, NewNoDbQB().Custom(customQ).Where(whereValues...))
	qb.Func(baseSQ).Merge()
	res, err := qb.Exec()
	if err != nil {
		return false, err
	}
	ok := res["rowsAffected"].(int64) != 0
	return ok, nil
}
