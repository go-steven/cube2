package main

import (
	"fmt"
	"github.com/go-steven/cube2/cube"
	"github.com/go-steven/cube2/engine"
	"github.com/go-steven/cube2/util"
)

func dimension_cube() cube.Cube {
	return cube.New().SQL(`SELECT 
		client_id, 
		campaign_id, 
		campaign_name, 
		record_on, 
		impressions, 
		uv, 
		click, 
		cost
	FROM skyline.zhizuan_campaign_rpt_daily
    WHERE client_id = 1
    	AND record_on BETWEEN '2017-03-06' AND '2017-03-12'
    ORDER BY record_on DESC
    LIMIT 0, 1
`)
}
func main() {
	reports := engine.NewReports()
	reports.AddCube("dimension example", dimension_cube())
	ret, err := reports.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("ret:", util.Json(ret))
}
