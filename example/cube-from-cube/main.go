package main

import (
	"fmt"
	"github.com/go-steven/cube2/cube"
	"github.com/go-steven/cube2/engine"
	"github.com/go-steven/cube2/util"
)

func cube_from_cube() cube.Cube {
	c := cube.New().SQL("SELECT * FROM skyline.clients")
	return cube.New().From(c)
}
func main() {
	reports := engine.NewReports()
	reports.AddCube("cube from cube example", cube_from_cube())
	ret, err := reports.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("ret:", util.Json(ret))
}
