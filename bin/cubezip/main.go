package main

import (
	"flag"
	"github.com/go-steven/cube2/util"
	"github.com/go-steven/cube2/util/dirzip"
	"github.com/go-steven/cube2/util/logger"
)

var (
	Logger     = logger.Logger
	logFlag    = flag.String("log", "", "set log path")
	dirFlag    = flag.String("dir", "", "set tpl script dir")
	outputFlag = flag.String("output", "", "set output file")
	mainFlag   = flag.Bool("main", true, "whether include the example main.go")
)

func main() {
	flag.Parse()
	Logger = logger.SetGlobalLogger(*logFlag)

	if *dirFlag == "" {
		flag.Usage()
		return
	}
	ret, err := dirzip.Zip(*dirFlag, *mainFlag)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Infof("Zip json: %s", ret)
	if *outputFlag != "" {
		if err := util.WriteFile(*outputFlag, []byte(ret)); err != nil {
			Logger.Error(err)
			return
		}
		Logger.Infof("Output zip json %s", *outputFlag)
	}
}
