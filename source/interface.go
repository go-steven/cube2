package source

import (
	"github.com/go-steven/cube2/util/logger"
	log "github.com/kdar/factorlog"
)

var (
	Logger *log.FactorLog = logger.Logger
)

func SetLogger(l *log.FactorLog) {
	Logger = l
}

type Row map[string]string
type Rows []Row

func (r Row) FieldsMapping(mapping map[string]string) Row {
	if mapping == nil || len(mapping) == 0 {
		return r
	}

	newRow := make(Row)
	for k, v := range r {
		newK, ok := mapping[k]
		if !ok {
			newK = k
		}
		newRow[newK] = v
	}
	return newRow
}

func (r Rows) FieldsMapping(mapping map[string]string) Rows {
	if mapping == nil || len(mapping) == 0 {
		return r
	}

	newRows := Rows{}
	for _, row := range r {
		newRows = append(newRows, row.FieldsMapping(mapping))
	}
	return newRows
}
