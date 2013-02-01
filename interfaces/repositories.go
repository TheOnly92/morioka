package interfaces

import (
	"github.com/TheOnly92/morioka/usecases"
	"errors"
	"github.com/gorilla/sessions"
	"net/http"
	"time"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

type DbProfile struct {
	Query     string
	Arguments []interface{}
	Duration  time.Duration
	Slow      bool
}

type RouteHandler struct {
	Path    string
	Handler func(w http.ResponseWriter, r *http.Request, user *usecases.User, title *Title) *WebError
	Methods []string
}

type TranslateHandler interface {
	GetSupportedLanguages() []string
	Get(key, lang string, a ...interface{}) string
	GetPlural(key, lang string, num int, a ...interface{}) string
}

type DbHandler interface {
	Begin() (Tx, error)
	Exec(query string, args ...interface{}) (Result, error)
	Query(query string, args ...interface{}) (Rows, error)
	QueryRow(query string, args ...interface{}) Row
	GetProfiledQueries() []DbProfile
}

type SessionHandler interface {
	GetSessionStore() sessions.Store
}

type DbRepo struct {
	db DbHandler
}

type Result interface {
	RowsAffected() (int64, error)
}

type Rows interface {
	Next() bool
	Err() error
	Scan(dest ...interface{}) error
}

type Row interface {
	Scan(dest ...interface{}) error
}

type Tx interface {
	Commit() error
	Exec(query string, args ...interface{}) (Result, error)
	Query(query string, args ...interface{}) (Rows, error)
	QueryRow(query string, args ...interface{}) Row
	Rollback() error
}
