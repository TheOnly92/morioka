package usecases

import (
	"bytes"
	"database/sql"
	"fmt"
	"runtime"
	"strings"
)

func IsCritical(err error) bool {
	if e, ok := err.(*DbError); ok {
		return !e.IsNotFound()
	}
	if v, ok := err.(Error); ok {
		if v.IsCritical() {
			return true
		}
	}
	return false
}

type Error interface {
	FileLine() (file string, line int)
	Trace() (stack []*runtime.Func)
	Package() string
	Function() string
	Tracef() string
	IsCritical() bool
	error
}

type errorBase struct {
	error
	*runtime.Func
	pc       []uintptr
	critical bool
}

func NewError(origErr error, critical bool) Error {
	if _, ok := origErr.(Error); ok {
		return origErr.(Error)
	}
	err := &errorBase{
		error:    origErr,
		pc:       make([]uintptr, 10),
		critical: critical,
	}

	var n int
	if n = runtime.Callers(2, err.pc); n > 0 {
		err.Func = runtime.FuncForPC(err.pc[0])
	}
	err.pc = err.pc[:n]

	return err
}

func (e *errorBase) IsCritical() bool {
	return e.critical
}

func (e *errorBase) Package() string {
	caller := strings.Split(e.Func.Name(), ".")
	return strings.Join(caller[0:len(caller)-1], ".")
}

func (e *errorBase) Function() string {
	caller := strings.Split(e.Func.Name(), ".")
	return caller[len(caller)-1]
}

func (e *errorBase) Error() string {
	return e.error.Error()
}

func (e *errorBase) FileLine() (file string, line int) {
	return e.Func.FileLine(e.pc[0])
}

func (e *errorBase) Trace() (stack []*runtime.Func) {
	stack = make([]*runtime.Func, len(e.pc))
	for i, pc := range e.pc {
		stack[i] = runtime.FuncForPC(pc)
	}

	return
}

func (e *errorBase) Tracef() string {
	depth := 10
	var last, name string
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "Trace: %s:\n", e.error.Error())
	for i, frame := range e.Trace() {
		if depth > 0 && i >= depth {
			break
		}
		file, line := frame.FileLine(e.pc[i])
		if name = frame.Name(); name != last {
			fmt.Fprintf(b, "\n %s:\n", frame.Name())
		}
		last = name
		fmt.Fprintf(b, "\t%s#L=%d\n", file, line)
	}

	return string(b.Bytes())
}

type DbError struct {
	*errorBase
	query      string
	isNotFound bool
}

func NewDbError(origErr error, query string) error {
	baseErr := &errorBase{
		pc:       make([]uintptr, 10),
		error:    origErr,
		critical: true,
	}

	var n int
	if n = runtime.Callers(2, baseErr.pc); n > 0 {
		baseErr.Func = runtime.FuncForPC(baseErr.pc[0])
	}
	baseErr.pc = baseErr.pc[:n]

	err := &DbError{
		errorBase:  baseErr,
		query:      query,
		isNotFound: (origErr == sql.ErrNoRows),
	}

	return err
}

func (e *DbError) IsCritical() bool {
	return !e.isNotFound
}

func (e *DbError) IsNotFound() bool {
	return e.isNotFound
}

func (e *DbError) Error() string {
	message := e.errorBase.Error()
	if e.query != "" {
		message += "\nQuery: " + e.query
	}
	return message
}
