package main

import (
	"github.com/go-steven/cube2/cube"
)

func zhizuan_daily_report() cube.Cube {
	c := cube.New().SQL(`SELECT 
			record_on,
			skyline.week_of_year(record_on) AS week_id,
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
		FROM skyline.zhizuan_campaign_rpt_daily 
		WHERE client_id = @CLIENT_ID@
			AND record_on BETWEEN '@START_DATE@' AND '@END_DATE@'
		GROUP BY record_on`).SQL(`SELECT
			record_on,
			CONCAT('第',week_id,'周') AS week_id,
			impressions,
			click,
			IF(impressions=0, 0, ROUND(100.0 * click / impressions, 2)) AS ctr,
			IF(click=0, 0, ROUND(1.0 * cost / click, 2)) AS cpc,
			ROUND(cost,2) AS cost,
			pay_count,
			ROUND(pay,2) AS pay,
			gmv_count,
			ROUND(gmv_amt,2) AS gmv_amt,
			cart,
			fav_shop_count,
			fav_item_count,
			IF(cost=0, 0, ROUND(1.0 * pay / cost, 2)) AS roi,
			IF(click=0, 0, ROUND(100.0 * pay_count / click, 2)) AS cvr,
			IF(pay_count=0, 0, ROUND(1.0 * cost / pay_count, 2)) AS cpu,
			IF(impressions=0, 0, ROUND(1000.0 * cost / impressions, 2)) AS cpm
		FROM @THIS@ AS t
	`)

	return c.RetMapping(map[string]string{
		"record_on":      "Date",
		"week_id":        "周数据",
		"impressions":    "展现",
		"click":          "点击量",
		"ctr":            "CTR",
		"cpc":            "CPC",
		"cost":           "花费",
		"pay_count":      "成交订单数",
		"pay":            "成交金额",
		"gmv_count":      "15订单数",
		"gmv_amt":        "15订单金额",
		"cart":           "添加购物车量",
		"fav_shop_count": "收藏店铺量",
		"fav_item_count": "收藏宝贝量",
		"roi":            "ROI",
		"cvr":            "CVR",
		"cpu":            "CPU",
		"cpm":            "CPM",
	})
}
