package main

import (
	"flag"
	"github.com/go-steven/cube2/engine"
	"github.com/go-steven/cube2/util"
	"github.com/go-steven/cube2/util/logger"
)

var (
	Logger      = logger.Logger
	logFlag     = flag.String("log", "", "set log path")
	zipjsonFlag = flag.String("zipjson", "", "set zip json file")
	tplCfgFlag  = flag.String("tplcfg", "", "set tpl cfg file")
	outputFlag  = flag.String("output", "", "set output file")
)

func main() {
	flag.Parse()
	if *zipjsonFlag == "" {
		flag.Usage()
		return
	}

	Logger = logger.SetGlobalLogger(*logFlag)
	engine.SetLogger(Logger)

	engine := engine.NewGoEngine()
	zipjson, err := util.ReadFile(*zipjsonFlag)
	if err != nil {
		Logger.Error(err)
		return
	}
	var tplcfg string
	if *tplCfgFlag != "" {
		content, err := util.ReadFile(*tplCfgFlag)
		if err != nil {
			Logger.Error(err)
			return
		}
		tplcfg = string(content)
	}
	ret, err := engine.Execute(string(zipjson), tplcfg)
	if err != nil {
		Logger.Error(err)
		return
	}
	Logger.Infof("ret: %s", util.Json(ret))
	if *outputFlag != "" {
		if err := util.WriteFile(*outputFlag, []byte(util.Json(ret))); err != nil {
			Logger.Error(err)
			return
		}
		Logger.Infof("Output reports json to %s", *outputFlag)
	}
}
