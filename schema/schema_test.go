package schema

import (
	"testing"

	"github.com/SeanChan0901/gee-orm/dialect"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

var (
	Sqlite3Dial, _ = dialect.GetDialect("sqlite3")
	MysqlDial, _   = dialect.GetDialect("mysql")
)

func TestParse(t *testing.T) {
	schema := Parse(&User{}, Sqlite3Dial)

	if schema.Name != "User" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}

	if schema.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
}
