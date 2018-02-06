package cube

import (
	"fmt"
	"github.com/bububa/mymysql/autorc"
	_ "github.com/bububa/mymysql/thrsafe"
	"github.com/go-steven/cube2/source"
	"github.com/go-steven/cube2/util"
	"github.com/go-steven/cube2/util/dbconn"
	"github.com/go-steven/cube2/util/errors"
	"strings"
	"sync"
)

type DefaultCube struct {
	db      *source.Mysql
	sql     string
	cubes   map[string]Cube
	summary map[string]Cube

	retMapping map[string]string

	m *sync.RWMutex // 用于并发
}

func New() *DefaultCube {
	return NewDefaultCube(dbconn.Mdb)
}

func NewDefaultCube(db *autorc.Conn) *DefaultCube {
	c := &DefaultCube{
		db:         source.NewMysql(db),
		cubes:      make(map[string]Cube),
		summary:    make(map[string]Cube),
		retMapping: make(map[string]string),

		m: new(sync.RWMutex),
	}
	c.Link(CubeTplVar(CUBE_THIS), c)
	return c
}

func (c *DefaultCube) FromTable(table string) Cube {
	if c == nil {
		return nil
	}
	c.set_cube_sql(fmt.Sprintf(`SELECT t.* FROM %s AS t`, table))
	return c
}

func (c *DefaultCube) From(cube Cube) Cube {
	if c == nil {
		return nil
	}

	c.set_cube_sql(cube.ToSQL())
	return c
}

func (c *DefaultCube) SQL(sql string, a ...interface{}) Cube {
	if c == nil {
		return nil
	}
	if len(a) > 0 {
		sql = fmt.Sprintf(sql, a...)
	}
	c.m.RLock()
	for tpl_var, cube := range c.cubes {
		if cube == nil {
			continue
		}
		sql = strings.Replace(sql, tpl_var, fmt.Sprintf("(%s)", cube.ToSQL()), -1)
	}
	c.m.RUnlock()

	c.set_cube_sql(sql)
	return c
}

func (c *DefaultCube) SummarySQL(name string, sql string, a ...interface{}) Cube {
	summary, ok := c.summary[name]
	if !ok {
		summary = New().From(c)
		c.summary[name] = summary
		summary.Link(CubeTplVar(CUBE_CUBE), c)
	}
	summary.SQL(sql, a...)
	summary.Link(CubeTplVar(CUBE_SUMMARY), summary)
	Logger.Infof(summary.ToSQL())
	return c
}

func (c *DefaultCube) GroupSummary(name, method string, fields []string) Cube {
	name = util.Trim(name)
	method = util.UpperTrim(method)
	if len(fields) == 0 || name == "" || method == "" {
		return c
	}

	sql := fmt.Sprintf(`SELECT 
`)
	cnt := len(fields)
	method = c.Escape(method)
	for i, v := range fields {
		new_v := c.Escape(v)
		sql += fmt.Sprintf(`%s(%s) AS %s`, method, new_v, new_v)
		if i < cnt-1 {
			sql += ","
		}
		sql += fmt.Sprintf(` 
`)
	}
	sql += `FROM @CUBE@ AS s`
	return c.SummarySQL(name, sql)
}

func (c *DefaultCube) ContrastSummary(name string, fields []string) Cube {
	name = util.Trim(name)
	if len(fields) == 0 || name == "" {
		return c
	}
	sql := fmt.Sprintf(`SELECT 
`)
	cnt := len(fields)
	for i, v := range fields {
		new_v := c.Escape(v)
		sql += fmt.Sprintf(`IF(a.%s IS NULL OR a.%s=0, 0, ROUND((b.%s - a.%s) / a.%s, 4)) AS %s`, new_v, new_v, new_v, new_v, new_v, new_v)
		if i < cnt-1 {
			sql += ","
		}
		sql += fmt.Sprintf(` 
`)
	}
	sql += `FROM (SELECT 1 AS ttt_id, t.* FROM @CUBE@ AS t LIMIT 0, 1) AS a
    INNER JOIN (SELECT 1 AS ttt_id, t.* FROM @CUBE@ AS t LIMIT 1, 1) AS b ON b.ttt_id = a.ttt_id`
	return c.SummarySQL(name, sql)
}

func (c *DefaultCube) RetMapping(mapping map[string]string) Cube {
	c.retMapping = mapping
	for _, v := range c.summary {
		v.RetMapping(mapping)
	}

	return c
}

func (c *DefaultCube) Link(alias string, cube Cube) Cube {
	if c == nil {
		return nil
	}

	alias = strings.TrimSpace(alias)
	alias = strings.ToUpper(alias)
	c.m.Lock()
	c.cubes[alias] = cube
	c.m.Unlock()

	return c
}

func (c *DefaultCube) Replace(tplcfg TplCfg) Cube {
	sql := c.ToSQL()
	for tpl_var, tpl_val := range tplcfg {
		tpl_var = CubeTplVar(tpl_var)
		if tpl_var != "" {
			sql = strings.Replace(sql, tpl_var, fmt.Sprintf("%v", tpl_val), -1)
		}
	}
	c.set_cube_sql(sql)

	c.m.RLock()
	for _, v := range c.summary {
		sql := v.ToSQL()
		for tpl_var, tpl_val := range tplcfg {
			tpl_var = CubeTplVar(tpl_var)
			if tpl_var != "" {
				sql = strings.Replace(sql, tpl_var, fmt.Sprintf("%v", tpl_val), -1)
			}
		}
		v.SQL(sql)
	}
	c.m.RUnlock()
	return c
}

func (c *DefaultCube) ToSQL() string {
	if c == nil {
		return ""
	}

	return c.get_cube_sql()
}

func (c *DefaultCube) Rows() (source.Rows, error) {
	if c == nil {
		return nil, errors.New("empty cube.")
	}

	sql := c.ToSQL()
	if strings.Contains(sql, TPL_SEP) {
		return nil, errors.New("SQL still has variables.")
	}
	rows, err := c.db.Query(sql)
	if err != nil {
		return nil, err
	}

	return rows.FieldsMapping(c.retMapping), nil
}

func (c *DefaultCube) Row() (source.Row, error) {
	if c == nil {
		return nil, errors.New("empty cube.")
	}
	sql := c.ToSQL()
	if strings.Contains(sql, TPL_SEP) {
		return nil, errors.New("SQL still has variables.")
	}
	row, err := c.db.QueryFirst(sql)
	if err != nil {
		return nil, err
	}

	return row.FieldsMapping(c.retMapping), nil
}

func (c *DefaultCube) Summary() (ret map[string]source.Row, err error) {
	if c == nil {
		return nil, errors.New("empty cube.")
	}
	ret = make(map[string]source.Row)
	for k, v := range c.summary {
		row, err := v.Row()
		if err != nil {
			return nil, err
		}
		ret[k] = row
	}
	return ret, nil
}

func (c *DefaultCube) Fields() ([]string, error) {
	if c == nil {
		return nil, errors.New("empty cube.")
	}
	sql := c.ToSQL()
	if strings.Contains(sql, TPL_SEP) {
		return nil, errors.New("SQL still has variables.")
	}
	fields, err := c.db.Fields(sql)
	if err != nil {
		return nil, err
	}
	ret := []string{}
	for _, v := range fields {
		new_v, ok := c.retMapping[v]
		if ok {
			ret = append(ret, new_v)
		} else {
			ret = append(ret, v)
		}
	}
	return ret, nil
}

func (c *DefaultCube) Escape(s string) string {
	return c.db.Escape(s)
}

func (c *DefaultCube) EscapeFields(fields []string) []string {
	ret := []string{}
	for _, v := range fields {
		ret = append(ret, c.db.EscapeFieldName(v))
	}

	return ret
}

func (c *DefaultCube) get_cube_sql() string {
	c.m.RLock()
	sql := c.sql
	c.m.RUnlock()

	return sql
}

func (c *DefaultCube) set_cube_sql(sql string) {
	c.m.Lock()
	c.sql = sql
	c.m.Unlock()
}

func (c *DefaultCube) Copy() Cube {
	cube := &DefaultCube{
		db:         c.db,
		cubes:      make(map[string]Cube),
		summary:    make(map[string]Cube),
		retMapping: make(map[string]string),
		m:          new(sync.RWMutex),
	}
	cube.Link(CubeTplVar(CUBE_THIS), cube)
	c.m.RLock()
	cube.sql = c.sql
	cube.retMapping = c.retMapping
	for k, v := range c.cubes {
		cube.cubes[k] = v
	}
	for k, v := range c.summary {
		v.Link(CubeTplVar(CUBE_CUBE), cube)
		v.Link(CubeTplVar(CUBE_SUMMARY), v)
		cube.summary[k] = v
	}
	c.m.RUnlock()
	return cube
}
