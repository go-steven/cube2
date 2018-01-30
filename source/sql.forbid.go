package source

import (
	"bytes"
	"fmt"
	"github.com/go-steven/cube2/util"
	"github.com/go-steven/cube2/util/errors"
	"regexp"
	"strings"
)

const (
	WITH_NO_SPACE = "no-space"
	WITH_SPACES   = "has_spaces"
)

var forbid_keywords = map[string][]string{
	WITH_NO_SPACE: []string{
		";",
	},
	WITH_SPACES: []string{
		"ALTER",
		"AUTOCOMMIT",
		"BEGIN",
		"CALL",
		"COMMIT",
		"CREATE",
		"DELETE",
		"DENY",
		"DO",
		"DROP",
		"END",
		"EXPLAIN",
		"FLUSH",
		"FUNCTION",
		"GRANT",
		"HANDLER",
		"INDEX",
		"INSERT",
		"INTO",
		"LOAD",
		"LOCK",
		"PROCEDURE",
		"READ",
		"RENAME",
		"REPLACE",
		"REVOKE",
		"ROLLBACK",
		"SAVEPOINT",
		"SET",
		"SHOW",
		"START",
		"TABLE",
		"TABLES",
		"TRANSACTION",
		"TRUNCATE",
		"UNLOCK",
		"UPDATE",
		"VALUES",
		"VIEW",
		"WITH",
	},
}

func SqlForbidRegexp() *regexp.Regexp {
	var buffer bytes.Buffer
	buffer.WriteString(`.*`)
	for keywords_type, keywords := range forbid_keywords {
		for k, v := range keywords {
			if k > 0 {
				buffer.WriteString(`|`)
			}
			switch keywords_type {
			case WITH_NO_SPACE:
				buffer.WriteString(util.UpperTrim(v))
			case WITH_SPACES:
				buffer.WriteString(fmt.Sprintf(`\\s%v\\s`, util.UpperTrim(v)))
			}
		}
	}
	buffer.WriteString(`.*`)
	//Logger.Infof("regexp: %v", buffer.String())
	reg, err := regexp.Compile(buffer.String())
	if err != nil {
		Logger.Error(errors.New(err.Error()))
		return nil
	}

	return reg
}

func SqlForbid(sql string, reg *regexp.Regexp) bool {
	sql = util.UpperTrim(sql)
	//Logger.Infof("SQL:%v", sql)
	if !strings.HasPrefix(sql, "SELECT") {
		Logger.Error(errors.New("SQL should begin with SELECT."))
		return true
	}

	if reg != nil && reg.MatchString(sql) {
		Logger.Error(errors.New("SQL has not-allowed keywords."))
		return true
	}

	return false
}
