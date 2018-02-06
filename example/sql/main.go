package main

import (
	"fmt"
	"github.com/go-steven/cube2/cube"
	"github.com/go-steven/cube2/engine"
	"github.com/go-steven/cube2/util"
)

func sql_cube() cube.Cube {
	return cube.New().SQL(`SELECT
		record_on, 
		SUM(impressions) AS impressions, 
		SUM(uv) AS uv, 
		SUM(click) AS click, 
		SUM(cost) AS cost, 
		SUM(cart) AS cart, 
		SUM(fav_item_count) AS fav_item_count, 
		SUM(fav_shop_count) AS fav_shop_count, 
		SUM(gmv_amt) AS gmv_amt, 
		SUM(gmv_count) AS gmv_count,
		SUM(pay) AS pay,
		SUM(pay_count) AS pay_count 
FROM skyline.zhizuan_client_rpt_daily  
WHERE client_id = 1 
	AND record_on between '2017-03-01' AND '2017-03-31' 
GROUP BY record_on 
ORDER BY record_on
`)
}

func main() {
	reports := engine.NewReports()
	reports.AddCube("SQL example", sql_cube())
	ret, err := reports.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("ret:", util.Json(ret))
}
