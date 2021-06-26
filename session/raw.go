package session

import (
	"database/sql"
	"strings"

	"github.com/SeanChan0901/geeorm/clause"
	"github.com/SeanChan0901/geeorm/dialect"
	"github.com/SeanChan0901/geeorm/schema"

	"github.com/SeanChan0901/geeorm/log"
)

type Session struct {
	db       *sql.DB
	tx       *sql.Tx
	dialect  dialect.Dialect
	refTable *schema.Schema
	clause   clause.Clause
	sql      strings.Builder
	sqlVars  []interface{}
}

// CommonDB is a minimal function set of db
type CommonDB interface {
	// Query returns one or more matched records(rows) searched from db
	Query(query string, args ...interface{}) (*sql.Rows, error)

	// QueryRow returns one matched record(row) searched from db
	QueryRow(query string, args ...interface{}) *sql.Row

	// Exec executes a SQL sentence
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil) // sql.db implements Query(), QueryRow() and Exec()
var _ CommonDB = (*sql.Tx)(nil) // sql.tx implements Query(), QueryRow() and Exec()

// DB returns tx if a tx begins. otherwise return *sql.DB
func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

// New returns a new sql session
func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

// Clear resets a session (reset sql sentence)
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause.Clear()
}

// Raw returns a session represent a raw sql sentence
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// Exec raw sql with sqlVars
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// QueryRow gets a record from db
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// QueryRows gets a list of record from db
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
