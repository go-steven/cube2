package engine

import (
	"github.com/go-steven/cube2/cube"
	"github.com/go-steven/cube2/source"
	"github.com/go-steven/cube2/util/logger"
	log "github.com/kdar/factorlog"
)

var (
	Logger *log.FactorLog = logger.Logger
)

var DEFAULT_LOG string = getDefaultTmpDir() + "cube.log"
/*
func init() {
	Logger = logger.SetGlobalLogger(DEFAULT_LOG)
	source.SetLogger(Logger)
	cube.SetLogger(Logger)
}
*/
func SetLogger(l *log.FactorLog) {
	Logger = l
	source.SetLogger(l)
	cube.SetLogger(l)
}
