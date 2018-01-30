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
	SQLCfg(tplcfg TplCfg) Cube
	SummarySQL(name string, sql string, a ...interface{}) Cube
	RetFieldsMapping(mapping map[string]string) Cube
	SummaryFieldsMapping(name string, mapping map[string]string) Cube
	Link(alias string, cube Cube) Cube

	ToSQL() string
	GetRows() (source.Rows, error)
	GetRow() (source.Row, error)
	Fields() ([]string, error)
	GetSummary() (map[string]source.Row, error)

	Escape(s string) string
	EscapeFields(fields []string) []string

	Copy() Cube
}

type Cuber func() Cube

func CubeTplVar(name string) string {
	name = util.UpperTrim(name)
	return fmt.Sprintf("%s%s%s", TPL_SEP, name, TPL_SEP)
}
