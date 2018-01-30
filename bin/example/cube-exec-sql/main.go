package main

import (
	"fmt"
	"github.com/go-steven/cube2/cube"
	"github.com/go-steven/cube2/engine"
	"github.com/go-steven/cube2/util"
)

func cube_exec_sql() cube.Cube {
	c := cube.New().SQL("SELECT * FROM skyline.clients")
	return c.SQL(`SELECT t.* FROM @THIS@ AS t WHERE t.id=10`)
}
func main() {
	reports := engine.NewReports()
	reports.AddCube("cube exec SQL example", cube_exec_sql())
	ret, err := reports.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("ret:", util.Json(ret))
}
