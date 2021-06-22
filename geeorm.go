package geeorm

import (
	"database/sql"
	"github.com/SeanChan0901/gee-orm/dialect"

	"github.com/SeanChan0901/gee-orm/log"
	"github.com/SeanChan0901/gee-orm/session"
	_ "github.com/mattn/go-sqlite3"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

type TxFunc func(*session.Session) (interface{}, error)

// Transaction starts a transaction.
// It does a Commit or Rollback depending on whether an error is returned.
func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()

	if err := s.Begin(); err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)  // re-throw panic after Rollback
		} else if err != nil {
			_ = s.Rollback()  // err is non-nil; don't change it
		} else {
			err = s.Commit()  // err is nil; if Commit returns error update err
		}
	}()

	return f(s)
}

// NewEngine returns a new engine instance connected to a given db
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}

	// Send a ping to make sure the database connection is alive
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	// make sure the specific dialect exists
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}

	e = &Engine{
		db:      db,
		dialect: dial,
	}
	log.Info("Connect database success")
	return
}

// Close closes the connection with db
func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

// NewSession returns a new db session
func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}
