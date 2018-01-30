package main

import (
	"fmt"
	"github.com/go-steven/cube2/engine"
	"github.com/go-steven/cube2/util"
)

func main() {
	curr_dir := util.CurrDir()
	fmt.Println("curr dir: ", curr_dir)

	script, err := util.ReadFile(fmt.Sprintf("%s/../report/mattel/main.go", curr_dir))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	tplcfg, err := util.ReadFile(fmt.Sprintf("%s/tpl.cfg", curr_dir))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	engine := engine.NewGoEngine()
	ret, err := engine.Execute(string(script), string(tplcfg))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("ret: ", util.Json(ret))
}
