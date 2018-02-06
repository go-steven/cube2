package cube

import (
	"fmt"
	"github.com/go-steven/cube2/source"
	"github.com/go-steven/cube2/util"
)

const (
	TPL_SEP      = "@"
	CUBE_THIS    = "THIS"
	CUBE_CUBE    = "CUBE"
	CUBE_SUMMARY = "SUMMARY"
)

type Cube interface {
	From(c Cube) Cube
	FromTable(table string) Cube
	SQL(sql string, a ...interface{}) Cube
	Replace(tplcfg TplCfg) Cube
	SummarySQL(name string, sql string, a ...interface{}) Cube
	GroupSummary(name, method string, fields []string) Cube
	ContrastSummary(name string, fields []string) Cube
	RetMapping(mapping map[string]string) Cube
	Link(alias string, cube Cube) Cube

	ToSQL() string
	Rows() (source.Rows, error)
	Row() (source.Row, error)
	Fields() ([]string, error)
	Summary() (map[string]source.Row, error)

	Escape(s string) string
	EscapeFields(fields []string) []string

	Copy() Cube
}

type Cuber func() Cube

func CubeTplVar(name string) string {
	name = util.UpperTrim(name)
	return fmt.Sprintf("%s%s%s", TPL_SEP, name, TPL_SEP)
}
