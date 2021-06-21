package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
}


func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i ++ {
		vars = append(vars, "?")
	}

	return strings.Join(vars, ",")
}

func _insert(values ...interface{}) (string, []interface{}) {
	// INSERT INTO $tableName ($fields)
	// e.g sql, _ := _insert("User", "Name", "Age")
	// SQL: "INSERT INTO User (Name, Age)"
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}

func _values(values ...interface{}) (string, []interface{}) {
	// VALUES（$v1), ($v2), ...
	// e.g sql, vars := _values("1", "2", "3")
	// SQL: "VALUES (?) (?) (?), "
	// vars: "1", "2", "3"
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		if bindStr  == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i + 1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

func _select(values ...interface{}) (string, []interface{}) {
	// SELECT $fields FROM $tableName
	// e.g sql, _ := _select("User", "Name", "Age")
	// SQL: "SELECT Name, Age FROM User"
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}

func _limit(values ...interface{}) (string, []interface{}) {
	// LIMIT $sum
	// e.g sql, vars := _limit(3)
	// SQL: "LIMIT ?"
	// vars: 3
	return "LIMIT ?", values
}

func _where(values ...interface{}) (string, []interface{}) {
	// WHERE $description
	// e.g sql, vars := _where("Name = ?", "Tom")
	// SQL: "WHERE Name = ?"
	// vars: "Tom"
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func _orderBy(values ...interface{}) (string, []interface{}) {
	// e.g sql, _ := _orderBy("Age ASC")
	// SQL: "ORDER BY Age ASC"
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}