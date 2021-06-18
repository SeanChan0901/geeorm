package session

import (
	"database/sql"
	"fmt"
	"github.com/SeanChan0901/gee-orm/dialect"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

var (
	TestDB      *sql.DB
	TestDial, _ = dialect.GetDialect("sqlite3")
)

func TestMain(m *testing.M) {
	var err error
	TestDB, err = sql.Open("sqlite3", "../gee.db")
	if err != nil {
		fmt.Println(err.Error())
	}
	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	return New(TestDB, TestDial)
}

func TestSession_Exec(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()

	if count, err := result.RowsAffected(); err != nil || count != 2 {
		t.Fatal("expect 2, but got", count)
	}
}

func TestSession_QueryRow(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	row := s.Raw("SELECT count(*) FROM User").QueryRow()
	var count int

	if err := row.Scan(&count); err != nil || count != 0 {
		t.Fatal("failed to query db", err)
	}
}

func TestSession_QueryRows(t *testing.T) {
	var names = [3]string{
		"Tom",
		"John",
		"Jack",
	}

	s := NewSession()
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
