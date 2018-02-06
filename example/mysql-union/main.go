package main

import (
	"fmt"
	"github.com/go-steven/cube2/cube"
	"github.com/go-steven/cube2/engine"
	"github.com/go-steven/cube2/util"
)

func mysql_union_cube() cube.Cube {
	return cube.New().SQL(`SELECT *
	FROM skyline.zhizuan_campaign_rpt_daily AS r1
    WHERE r1.client_id = 1
    	AND r1.record_on BETWEEN '2017-03-06' AND '2017-03-12'
    LIMIT 0, 1
   	UNION
	SELECT r2.*
	FROM skyline.zhizuan_campaign_rpt_daily AS r2
    WHERE r2.client_id = 1
    	AND r2.record_on BETWEEN '2017-03-06' AND '2017-03-12'
    LIMIT 0, 1
`)
}
func main() {
	reports := engine.NewReports()
	reports.AddCube("mysql union example", mysql_union_cube())
	ret, err := reports.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("ret:", util.Json(ret))
}
