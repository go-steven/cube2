package main

import (
	"fmt"
	"github.com/go-steven/cube2/cube"
	"github.com/go-steven/cube2/engine"
	"github.com/go-steven/cube2/util"
)

func cube_from_sql() cube.Cube {
	return cube.New().SQL("SELECT * FROM skyline.clients")
}
func main() {
	reports := engine.NewReports()
	reports.AddCube("cube from SQL example", cube_from_sql())
	ret, err := reports.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("ret:", util.Json(ret))
}
