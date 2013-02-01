package infrastructure

import (
	"database/sql"
	"github.com/TheOnly92/morioka/interfaces"
	"github.com/TheOnly92/morioka/usecases"
	_ "github.com/bmizerany/pq"
	"log"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var Profile = false

var ProfiledQueries []interfaces.DbProfile

type PgsqlHandler struct {
	Conn *sql.DB
}

type PgsqlResult struct {
	res sql.Result
}

type PgsqlRows struct {
	rows *sql.Rows
}

type PgsqlRow struct {
	row *sql.Row
}

type PgsqlTx struct {
	tx *sql.Tx
}

func (handler PgsqlHandler) GetProfiledQueries() []interfaces.DbProfile {
	if !Profile {
		return []interfaces.DbProfile{}
	}
	return ProfiledQueries
}

func (handler *PgsqlHandler) Exec(query string, args ...interface{}) (interfaces.Result, error) {
	var start, end int64
	queryStr, passArgs := processQueryAndArguments(query, args)
	if Profile {
		start = time.Now().UnixNano()
	}
	result, err := handler.Conn.Exec(queryStr, passArgs...)
	if Profile {
		end = time.Now().UnixNano()
		ProfiledQueries = append(ProfiledQueries, interfaces.DbProfile{queryStr, passArgs, time.Duration(end-start) * time.Nanosecond, time.Duration(end-start)*time.Nanosecond > 500*time.Millisecond})
	}
	if err != nil {
		err = usecases.NewDbError(err, query)
		LogSQLError(err)
	}
	return result, err
}

func (handler *PgsqlHandler) Query(query string, args ...interface{}) (interfaces.Rows, error) {
	var start, end int64
	queryStr, passArgs := processQueryAndArguments(query, args)
	if Profile {
		start = time.Now().UnixNano()
	}
	rows, err := handler.Conn.Query(queryStr, passArgs...)
	if Profile {
		end = time.Now().UnixNano()
		ProfiledQueries = append(ProfiledQueries, interfaces.DbProfile{queryStr, passArgs, time.Duration(end-start) * time.Nanosecond, time.Duration(end-start)*time.Nanosecond > 500*time.Millisecond})
	}
	if err != nil {
		err = usecases.NewDbError(err, query)
		LogSQLError(err)
	}
	return &PgsqlRows{rows}, err
}

func (handler *PgsqlHandler) QueryRow(query string, args ...interface{}) interfaces.Row {
	var start, end int64
	queryStr, passArgs := processQueryAndArguments(query, args)
	if Profile {
		start = time.Now().UnixNano()
	}
	row := handler.Conn.QueryRow(queryStr, passArgs...)
	if Profile {
		end = time.Now().UnixNano()
		ProfiledQueries = append(ProfiledQueries, interfaces.DbProfile{queryStr, passArgs, time.Duration(end-start) * time.Nanosecond, time.Duration(end-start)*time.Nanosecond > 500*time.Millisecond})
	}
	return &PgsqlRow{row}
}

func (handler *PgsqlHandler) Begin() (interfaces.Tx, error) {
	tx, err := handler.Conn.Begin()
	if err != nil {
		err = usecases.NewDbError(err, "")
		LogSQLError(err)
	}
	return &PgsqlTx{tx}, err
}

func processQueryAndArguments(query string, args []interface{}) (string, []interface{}) {
	var passArgs []interface{}
	j := 1
	for i, v := range args {
		t := reflect.ValueOf(v)
		switch t.Kind() {
		default:
			passArgs = append(passArgs, v)
			reg := regexp.MustCompile("\\$" + strconv.Itoa(i+1) + "([^0-9]|$)")
			query = reg.ReplaceAllString(query, "^^^"+strconv.Itoa(j)+"$1")
			j++
		case reflect.Slice:
			var replace []string
			for _, v2 := range v.([]interface{}) {
				passArgs = append(passArgs, v2)
				replace = append(replace, "^^^"+strconv.Itoa(j))
				j++
			}
			reg := regexp.MustCompile("\\$" + strconv.Itoa(i+1) + "([^0-9])")
			query = reg.ReplaceAllString(query, strings.Join(replace, ",")+"$1")
		}
	}
	query = strings.Replace(query, "^^^", "$", -1)
	return query, passArgs
}

func (r PgsqlRows) Next() bool {
	return r.rows.Next()
}

func (r PgsqlRows) Err() error {
	err := r.rows.Err()
	if err != nil {
		err = usecases.NewDbError(err, "")
		LogSQLError(err)
	}
	return err
}

func (r PgsqlRows) Scan(dest ...interface{}) error {
	err := r.rows.Scan(dest...)
	if err != nil {
		err = usecases.NewDbError(err, "")
		LogSQLError(err)
	}
	return err
}

func (r PgsqlRow) Scan(dest ...interface{}) error {
	err := r.row.Scan(dest...)
	if err != nil {
		err = usecases.NewDbError(err, "")
		LogSQLError(err)
	}
	return err
}

func (t PgsqlTx) Commit() error {
	err := t.tx.Commit()
	if err != nil {
		err = usecases.NewDbError(err, "")
		LogSQLError(err)
	}
	return err
}

func (t PgsqlTx) Exec(query string, args ...interface{}) (interfaces.Result, error) {
	queryStr, passArgs := processQueryAndArguments(query, args)
	ret, err := t.tx.Exec(queryStr, passArgs...)
	if err != nil {
		err = usecases.NewDbError(err, query)
		LogSQLError(err)
	}
	return &PgsqlResult{ret}, err
}

func (t PgsqlTx) Query(query string, args ...interface{}) (interfaces.Rows, error) {
	queryStr, passArgs := processQueryAndArguments(query, args)
	rows, err := t.tx.Query(queryStr, passArgs...)
	if err != nil {
		err = usecases.NewDbError(err, query)
		LogSQLError(err)
	}
	return &PgsqlRows{rows}, err
}

func (t PgsqlTx) QueryRow(query string, args ...interface{}) interfaces.Row {
	queryStr, passArgs := processQueryAndArguments(query, args)
	return &PgsqlRow{t.tx.QueryRow(queryStr, passArgs...)}
}

func (t PgsqlTx) Rollback() error {
	err := t.tx.Rollback()
	if err != nil {
		err = usecases.NewDbError(err, "")
		LogSQLError(err)
	}
	return err
}

func (r PgsqlResult) RowsAffected() (int64, error) {
	rt, err := r.res.RowsAffected()
	if err != nil {
		err = usecases.NewDbError(err, "")
		LogSQLError(err)
	}
	return rt, err
}

func NewPgsqlHandler(pgsqlUrl string) *PgsqlHandler {
	pgsql, _ := url.Parse(pgsqlUrl)
	dbuser := pgsql.User.Username()
	password, _ := pgsql.User.Password()
	host := strings.Split(pgsql.Host, ":")
	dbname := strings.TrimLeft(pgsql.Path, "/")
	db, err := sql.Open("postgres", "user="+dbuser+" dbname="+dbname+" password="+password+" sslmode=disable port="+host[1]+" host="+host[0])
	if err != nil {
		log.Fatal("Connect DB ", err)
	}
	return &PgsqlHandler{db}
}

func LogSQLError(err error) {}
