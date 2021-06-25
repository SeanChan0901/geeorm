package session

import "testing"

type User struct {
	Name string `geerom:"PRIMARY KEY"`
	Age  int
}

func testMysqlTableInit(t *testing.T) *Session {
	t.Helper()
	s := NewMysqlSession().Model(&User{})
	return s
}

func testSqliteTableInit(t *testing.T) *Session {
	t.Helper()
	s := NewSqlite3Session().Model(&User{})
	return s
}

func TestSession_CreateTable(t *testing.T) {
	sqlite3Session := testSqliteTableInit(t) // 创建一个关联 User 表的数据库 Session
	_ = sqlite3Session.DropTable()
	_ = sqlite3Session.CreateTable()
	if !sqlite3Session.HasTable() {
		t.Fatal("Failed to create sqlite3 table User")
	}

	mysqlSession := testMysqlTableInit(t)
	_ = mysqlSession.DropTable()
	_ = mysqlSession.CreateTable()
	if !mysqlSession.HasTable() {
		t.Fatal("Failed to create mysql table User")
	}
}
