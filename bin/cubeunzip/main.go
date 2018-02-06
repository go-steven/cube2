package main

import (
	"flag"
	"github.com/go-steven/cube2/util"
	"github.com/go-steven/cube2/util/dirzip"
	"github.com/go-steven/cube2/util/logger"
)

var (
	Logger      = logger.Logger
	logFlag     = flag.String("log", "", "set log path")
	zipjsonFlag = flag.String("zipjson", "", "set zip json for tpl script dir")
	outputFlag  = flag.String("output", "", "set output dir")
	mainFlag    = flag.Bool("main", true, "whether include the example main.go")
	dirname     string
)

func main() {
	flag.Parse()
	Logger = logger.SetGlobalLogger(*logFlag)

	zip, err := util.ReadFile(*zipjsonFlag)
	if err != nil {
		panic(err)
	}
	if err = dirzip.UnZip(string(zip), *outputFlag, *mainFlag); err != nil {
		panic(err)
	}
	Logger.Infof("Unzipped json zip to dir: %s", *outputFlag)
}
