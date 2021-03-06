package main

import (
	"fmt"
	"github.com/go-steven/cube2/cube"
	"github.com/go-steven/cube2/util"
)

// when testing, rename it to main()
func main() {
	tplCfgs := cube.TplCfg{
		"START_DATE": "2017-12-01",
		"END_DATE":   "2017-12-31",
		"CLIENT_ID":  10,
	}
	r := Reports()
	ret, err := r.RunWithCfgs(tplCfgs)
	if err != nil {
		panic(err)
	}
	fmt.Println("ret:", util.Json(ret))
}
