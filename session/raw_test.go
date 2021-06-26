package session

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/SeanChan0901/geeorm/dialect"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var (
	TestSqlit3DB   *sql.DB
	TestMysqlDB    *sql.DB
	Sqlite3Dial, _ = dialect.GetDialect("sqlite3")
	MysqlDial, _   = dialect.GetDialect("mysql")
)

func TestMain(m *testing.M) {
	var err error
	TestSqlit3DB, err = sql.Open("sqlite3", "../gee.db")
	if err != nil {
		fmt.Println(err.Error())
	}

	TestMysqlDB, err = sql.Open("mysql", "sean:123456@tcp(127.0.0.1:3306)/geeorm_test")
	if err != nil {
		fmt.Println(err.Error())
	}

	code := m.Run()
	_ = TestSqlit3DB.Close()
	os.Exit(code)
}

func NewSqlite3Session() *Session {
	return New(TestSqlit3DB, Sqlite3Dial)
}

func NewMysqlSession() *Session {
	return New(TestMysqlDB, MysqlDial)
}

func TestSqlite3Session_Exec(t *testing.T) {
	s := NewSqlite3Session()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()

	if count, err := result.RowsAffected(); err != nil || count != 2 {
		t.Fatal("expect 2, but got", count)
	}
}

func TestMysqlSession_Exec(t *testing.T) {
	s := NewMysqlSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()

	if count, err := result.RowsAffected(); err != nil || count != 2 {
		t.Fatal("expect 2, but got", count)
	}
}

func TestSqlite3Session_QueryRow(t *testing.T) {
	s := NewSqlite3Session()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	row := s.Raw("SELECT count(*) FROM User").QueryRow()
	var count int

	if err := row.Scan(&count); err != nil || count != 0 {
		t.Fatal("failed to query db", err)
	}
}

func TestMysqlSession_QueryRow(t *testing.T) {
	s := NewMysqlSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	row := s.Raw("SELECT count(*) FROM User").QueryRow()
	var count int

	if err := row.Scan(&count); err != nil || count != 0 {
		t.Fatal("failed to query db", err)
	}
}

func TestSqlite3Session_QueryRows(t *testing.T) {
	var names = [3]string{
		"Tom",
		"John",
		"Jack",
	}

	s := NewSqlite3Session()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("INSERT INTO User(`Name`) values (?), (?), (?)", names[0], names[1], names[2]).Exec()
	rows, err := s.Raw("SELECT * FROM User").QueryRows()

	if err != nil {
		t.Fatal("failed to query db", err)
	}
	var name string
	counter := 0

	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			t.Fatal("failed to scan rows", err)
		}
		if name != names[counter] {
			err = fmt.Errorf("expect %s, but got %s", names[counter], name)
			t.Fatal("data is not correct", err)
		}
		counter++
	}
}

func TestMysqlSession_QueryRows(t *testing.T) {
	var names = [3]string{
		"Tom",
		"John",
		"Jack",
	}

	s := NewMysqlSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("INSERT INTO User(`Name`) values (?), (?), (?)", names[0], names[1], names[2]).Exec()
	rows, err := s.Raw("SELECT * FROM User").QueryRows()

	if err != nil {
		t.Fatal("failed to query db", err)
	}
	var name string
	counter := 0

	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			t.Fatal("failed to scan rows", err)
		}
		if name != names[counter] {
			err = fmt.Errorf("expect %s, but got %s", names[counter], name)
			t.Fatal("data is not correct", err)
		}
		counter++
	}
}
