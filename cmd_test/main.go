package main

import (
	"database/sql"
	"fmt"

	geeorm "github.com/SeanChan0901/gee-orm"
)

func main() {
	db, _ := sql.Open("sqlite3", "fee.db")
	_, _ = db.Exec("")
	engine, _ := geeorm.NewEngine("sqlite3", "gee.db")
	defer engine.Close()

	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	result, _ := s.Raw("INSERT INTO User('Name') values (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}
