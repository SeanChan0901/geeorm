package clause

import (
	"reflect"
	"testing"
)

func testSelect(t *testing.T) {
	var clause Clause
	clause.Set(LIMIT, 3)
	clause.Set(SELECT, "User", []string{"Name", "Age"})
	clause.Set(WHERE, "Name = ?", "Tom")
	clause.Set(GROUPBY, "Age")
	clause.Set(ORDERBY, "Age ASC")
	sql, vars := clause.Build(SELECT, WHERE, GROUPBY, ORDERBY, LIMIT)
	t.Log(sql, vars)

	if sql != "SELECT Name,Age FROM User WHERE Name = ? GROUP BY Age ORDER BY Age ASC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}

	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}

func TestClause_Build(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		testSelect(t)
	})
}
