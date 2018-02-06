package main

import (
	"fmt"
	"github.com/go-steven/cube2/cube"
	"github.com/go-steven/cube2/engine"
	"github.com/go-steven/cube2/util"
)

func summary_cube() cube.Cube {
	c := cube.New().SQL(`SELECT
		record_on,
		SUM(impressions) AS total_impressions,
		SUM(uv) AS total_uv,
		SUM(click) AS total_click,
		SUM(cost) AS total_cost
	FROM skyline.zhizuan_campaign_rpt_daily
    WHERE client_id = 1
    	AND record_on BETWEEN '2017-03-06' AND '2017-03-12'
    GROUP BY record_on
`)

	c.SummarySQL("总计", `SELECT SUM(total_impressions) AS total_impressions,
		SUM(total_uv) AS total_uv,
		SUM(total_click) AS total_click,
		SUM(total_cost) AS total_cost,
		ROUND(IF(SUM(total_impressions)=0, 0, SUM(total_click) / SUM(total_impressions)),2) AS cpc
FROM @CUBE@ AS t`)
	return c
}
func main() {
	reports := engine.NewReports()
	reports.AddCube("summary example", summary_cube())
	ret, err := reports.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("ret:", util.Json(ret))
}
