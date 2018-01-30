package engine

import (
	"github.com/go-steven/cube2/source"
	"github.com/go-steven/cube2/util/logger"
	log "github.com/kdar/factorlog"
)

var (
	Logger *log.FactorLog = logger.Logger
)

func SetLogger(l *log.FactorLog) {
	Logger = l
	source.SetLogger(l)
}
