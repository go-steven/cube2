package errors

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

const (
	DEFAULT_ERR_CODE   = 1
	DEFAULT_CALL_DEPTH = 1
)

type CErr struct {
	File string
	Line int

	Code int
	Msg  string
}

func newCErr(code int, msg string, calldepth int) *CErr {
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}

	return &CErr{
		File: skip_gopath(file),
		Line: line,
		Code: code,
		Msg:  msg,
	}
}

func (e *CErr) Error() string {
	if e.Code == DEFAULT_ERR_CODE {
		return fmt.Sprintf("[%s:%d]%s", e.File, e.Line, e.Msg)
	} else {
		return fmt.Sprintf("[%s:%d]%d:%s", e.File, e.Line, e.Code, e.Msg)
	}
}

func is_windows() bool {
	return runtime.GOOS == "windows"
}

func skip_gopath(file string) string {
	isWindows := is_windows()
	sep := ":"
	if isWindows {
		sep = ";"
	}
	dirs := strings.Split(os.Getenv("GOPATH"), sep)
	for _, v := range dirs {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if isWindows {
			v = strings.Replace(v, "\\", "/", -1)
		}
		v += "/src/"
		if strings.HasPrefix(file, v) {
			return file[len(v):]
		}
	}

	return file
}

func New(text string) error {
	return newCErr(DEFAULT_ERR_CODE, text, DEFAULT_CALL_DEPTH+1)
}

func NewErr(err error) error {
	return newCErr(DEFAULT_ERR_CODE, err.Error(), DEFAULT_CALL_DEPTH+1)
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
func Errorf(format string, a ...interface{}) error {
	return newCErr(DEFAULT_ERR_CODE, fmt.Sprintf(format, a...), DEFAULT_CALL_DEPTH+1)
}

func NewCode(code int, text string) error {
	return newCErr(code, text, DEFAULT_CALL_DEPTH+1)
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
func ErrorfCode(code int, format string, a ...interface{}) error {
	return newCErr(code, fmt.Sprintf(format, a...), DEFAULT_CALL_DEPTH+1)
}
