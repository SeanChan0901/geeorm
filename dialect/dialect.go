package dialect

import "reflect"

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	// DataTypeOf returns database data type from go data type.
	// e.g reflect.Int in golang means "Integer" in sqlite3.
	DataTypeOf(typ reflect.Value) string

	// TableExistSQL returns a SQL sentence.
	// The SQL sentence it return can be used to check a if a table exists.
	TableExistSQL(tableName string) (string, []interface{})
}

// RegisterDialect add a new dialect to the known dialectsMap.
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// GetDialect returns a named dialect if it exists.
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
