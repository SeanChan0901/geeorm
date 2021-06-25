package geeorm

import (
	"errors"
	"reflect"
	"testing"

	"github.com/SeanChan0901/gee-orm/session"
)

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func OpenSqlite3DB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "gee.db")
	if err != nil {
		t.Fatal("failed to connect sqlite3", err)
	}
	return engine
}

func OpenMysqlDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("mysql", "sean:123456@tcp(127.0.0.1:3306)/geeorm_test")
	if err != nil {
		t.Fatal("failed to connect mysql", err)
	}
	return engine
}

func TestNewEngine(t *testing.T) {
	sqlite3Engine := OpenSqlite3DB(t)
	mysqlEngine := OpenMysqlDB(t)
	defer func() {
		sqlite3Engine.Close()
		mysqlEngine.Close()
	}()
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("sqlite3_rollback", func(t *testing.T) {
		sqlite3TransactionRollback(t)
	})
	t.Run("sqlite3_commit", func(t *testing.T) {
		sqlite3TransactionCommit(t)
	})
	t.Run("mysql_rollback", func(t *testing.T) {
		mysqlTransactionRollback(t)
	})
	t.Run("mysql_commit", func(t *testing.T) {
		mysqlTransactionCommit(t)
	})
}

func sqlite3TransactionRollback(t *testing.T) {
	engine := OpenSqlite3DB(t)
	defer engine.Close()

	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()

	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return nil, errors.New("error occurs") // returns a error to test rollback
	})

	if err == nil || s.HasTable() {
		t.Fatal("failed to rollback")
	}
}

func sqlite3TransactionCommit(t *testing.T) {
	engine := OpenSqlite3DB(t)
	defer engine.Close()

	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return
	})

	u := &User{}
	_ = s.First(u)
	if err != nil || u.Name != "Tom" {
		t.Fatal("failed to commit")
	}
}

func mysqlTransactionRollback(t *testing.T) {
	engine := OpenMysqlDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_ = s.Model(&User{}).CreateTable()
	// create/drop tables in mysql will instantly commit,
	// so create/drop operations do not support rollback.
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{})
		_, err = s.Insert(&User{"Tom", 18})
		_, err = s.Insert(&User{"John", 18})
		_, err = s.Insert(&User{"Tony", 18})
		_, err = s.Insert(&User{"Sean", 18})
		return nil, errors.New("Error")
	})

	if err == nil {
		t.Fatal("failed to rollback, expect err")
	}

	if count, err := s.Count(); err != nil && count != 0 {
		t.Fatal("failed to rollback, expect count = 0, but got", count)
	}
}

func mysqlTransactionCommit(t *testing.T) {
	engine := OpenMysqlDB(t)
	defer engine.Close()

	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return
	})

	u := &User{}
	_ = s.First(u)
	if err != nil || u.Name != "Tom" {
		t.Fatal("failed to commit")
	}
}

func TestSqlite3Engine_Migrate(t *testing.T) {
	engine := OpenSqlite3DB(t)
	defer engine.Close()

	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text PRIMARY KEY, XXX integer);").Exec()
	_, _ = s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	engine.Migrate(&User{})

	rows, _ := s.Raw("SELECT * FROM User").QueryRows()
	columns, _ := rows.Columns()
	if !reflect.DeepEqual(columns, []string{"Name", "Age"}) {
		t.Fatal("Failed to migrate table User, got columns", columns)
	}
}

func TestMysqlEngine_Migrate(t *testing.T) {
	engine := OpenMysqlDB(t)
	defer engine.Close()

	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name varchar(255) PRIMARY KEY, XXX integer);").Exec()
	_, _ = s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	engine.Migrate(&User{})

	rows, _ := s.Raw("SELECT * FROM User").QueryRows()
	columns, _ := rows.Columns()

	if !reflect.DeepEqual(columns, []string{"Name", "Age"}) {
		t.Fatal("Failed to migrate table User, got columns", columns)
	}
}
