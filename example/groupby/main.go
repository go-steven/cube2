package main

import (
	"fmt"
	"github.com/go-steven/cube2/cube"
	"github.com/go-steven/cube2/engine"
	"github.com/go-steven/cube2/util"
)

func groupby_cube() cube.Cube {
	return cube.New().SQL(`SELECT 
		campaign_id, 
		campaign_name, 
		record_on,
		SUM(impressions) AS total_impressions,
		SUM(uv) AS total_uv,
		SUM(click) AS total_click,
		SUM(cost) AS total_cost,
		AVG(impressions) AS avg_impressions,
		AVG(uv) AS avg_uv,
		AVG(click) AS avg_click,
		AVG(cost) AS avg_cost,
		MAX(impressions) AS max_impressions,
		MAX(uv) AS max_uv,
		MAX(click) AS max_click,
		MAX(cost) AS max_cost,
		MIN(impressions) AS min_impressions,
		MIN(uv) AS min_uv,
		MIN(click) AS min_click,
		MIN(cost) AS min_cost,
		COUNT(*) AS total_cnt
	FROM skyline.zhizuan_campaign_rpt_daily
    WHERE client_id = 1
    	AND record_on BETWEEN '2017-03-06' AND '2017-03-12'
    GROUP BY campaign_id, campaign_name, record_on
`)
}
func main() {
	reports := engine.NewReports()
	reports.AddCube("groupby example", groupby_cube())
	ret, err := reports.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("ret:", util.Json(ret))
}
