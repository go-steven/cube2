package source

import (
	"fmt"
	"github.com/bububa/mymysql/autorc"
	_ "github.com/bububa/mymysql/thrsafe"
	"github.com/go-steven/cube2/util/errors"
	"regexp"
)

type Mysql struct {
	db        *autorc.Conn
	sqlForbid *regexp.Regexp
}

func NewMysql(db *autorc.Conn) *Mysql {
	return &Mysql{
		db:        db,
		sqlForbid: SqlForbidRegexp(),
	}
}
func (m *Mysql) Query(sql string) (Rows, error) {
	//Logger.Infof("MYSQL run SQL: %s\n", sql)
	if m.sqlForbid != nil && SqlForbid(sql, m.sqlForbid) {
		Logger.Errorf("MYSQL run SQL: %s\n", sql)
		return nil, errors.New("ERROR: has not-allowed keywords in SQL.")
	}
	rows, res, err := m.db.Query(sql)
	if err != nil {
		return nil, errors.NewErr(err)
	}
	fields := []string{}
	for _, v := range res.Fields() {
		fields = append(fields, v.Name)
	}
	//Logger.Infof("fields: %s", util.Json(fields))
	if len(fields) == 0 {
		Logger.Errorf("MYSQL run SQL: %s\n", sql)
		return nil, errors.New("ERROR: no return fields given.")
	}

	for k, v := range fields {
		if v == "" {
			continue
		}
		if v[0] == '`' && v[len(v)-1] == '`' {
			v = v[1 : len(v)-1]
		}
		fields[k] = v
	}
	//Logger.Infof("new fields: %s", util.Json(fields))

	ret := Rows{}
	for _, row := range rows {
		retRow := make(Row)
		for _, v := range fields {
			retRow[v] = row.Str(res.Map(v))
		}

		ret = append(ret, retRow)
	}
	//Logger.Infof("SQL result: %d", len(ret))
	return ret, nil
}

func (m *Mysql) QueryFirst(sql string) (Row, error) {
	//Logger.Infof("MYSQL run SQL: %s\n", sql)
	if m.sqlForbid != nil && SqlForbid(sql, m.sqlForbid) {
		Logger.Errorf("MYSQL run SQL: %s\n", sql)
		return nil, errors.New("ERROR: has not-allowed keywords in SQL.")
	}
	row, res, err := m.db.QueryFirst(sql)
	if err != nil {
		return nil, errors.NewErr(err)
	}
	fields := []string{}
	for _, v := range res.Fields() {
		fields = append(fields, v.Name)
	}
	//Logger.Infof("fields: %s", util.Json(fields))
	if len(fields) == 0 {
		Logger.Infof("MYSQL run SQL: %s\n", sql)
		return nil, errors.New("ERROR: no return fields given.")
	}

	for k, v := range fields {
		if v == "" {
			continue
		}
		if v[0] == '`' && v[len(v)-1] == '`' {
			v = v[1 : len(v)-1]
		}
		fields[k] = v
	}
	//Logger.Infof("new fields: %s", util.Json(fields))

	ret := make(Row)
	for _, v := range fields {
		ret[v] = row.Str(res.Map(v))
	}
	//Logger.Infof("SQL result: %d", len(ret))
	return ret, nil
}

func (m *Mysql) Fields(sql string) ([]string, error) {
	//Logger.Infof("MYSQL run SQL: %s\n", sql)
	if m.sqlForbid != nil && SqlForbid(sql, m.sqlForbid) {
		Logger.Errorf("MYSQL run SQL: %s\n", sql)
		return nil, errors.New("ERROR: has not-allowed keywords in SQL.")
	}
	_, res, err := m.db.QueryFirst(sql)
	if err != nil {
		Logger.Errorf("MYSQL run SQL: %s\n", sql)
		return nil, errors.NewErr(err)
	}
	fields := []string{}
	for _, v := range res.Fields() {
		fields = append(fields, v.Name)
	}
	//Logger.Infof("fields: %s", util.Json(fields))
	if len(fields) == 0 {
		Logger.Errorf("MYSQL run SQL: %s\n", sql)
		return nil, errors.New("ERROR: no return fields given.")
	}

	for k, v := range fields {
		if v == "" {
			continue
		}
		if v[0] == '`' && v[len(v)-1] == '`' {
			v = v[1 : len(v)-1]
		}
		fields[k] = v
	}
	//Logger.Infof("new fields: %s", util.Json(fields))

	return fields, nil
}

func (m *Mysql) Escape(s string) string {
	return m.db.Escape(s)
}

func (m *Mysql) EscapeFieldName(s string) string {
	s = m.Escape(s)

	return fmt.Sprintf("`%s`", s)
}
