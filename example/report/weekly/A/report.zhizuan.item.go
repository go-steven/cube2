package main

import (
	"github.com/go-steven/cube2/cube"
)

// 钻展子分类表现
func zhizuan_item() cube.Cube {
	c := cube.New().SQL(`SELECT 
			CASE campaign_name
				WHEN REPLACE(campaign_name, " ", "") REGEXP ".*吸尘器.*" THEN "吸尘器"
				WHEN REPLACE(campaign_name, " ", "") REGEXP ".*除螨.*" THEN "除螨"
				WHEN REPLACE(campaign_name, " ", "") REGEXP ".*车载.*" THEN "车载"
				WHEN REPLACE(campaign_name, " ", "") REGEXP ".*吹风机.*" THEN "吹风机"
				WHEN REPLACE(campaign_name, " ", "") REGEXP ".*净化.*" THEN "净化"
				WHEN REPLACE(campaign_name, " ", "") REGEXP ".*圆筒.*" THEN "圆筒"
				WHEN REPLACE(campaign_name, " ", "") REGEXP ".*取暖.*" THEN "取暖"
				WHEN REPLACE(campaign_name, " ", "") REGEXP ".*全.*" THEN "全"
           ELSE ""
			END AS item,
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
			AND record_on BETWEEN '@START_DATE@' AND DATE_SUB(DATE_ADD('@START_DATE@', INTERVAL 1 @PERIOD@), INTERVAL 1 DAY)
		GROUP BY item`).SQL(`SELECT
			"钻展" AS platform,
			item,
			impressions,
			click,
			CASE WHEN impressions>0 THEN ROUND(click / impressions, 4) ELSE 0 END AS ctr,
			ROUND(cost,2) AS cost,
			CASE WHEN click>0 THEN ROUND(cost / click, 4) ELSE 0 END AS cpc,
			pay_count,
			ROUND(pay,2) AS pay,
			cart,
			fav_shop_count,
			CASE WHEN cost>0 THEN ROUND(pay / cost, 4) ELSE 0 END AS roi,
			CASE WHEN click>0 THEN ROUND(pay_count / click, 4) ELSE 0 END AS cvr,
			CASE WHEN pay_count>0 THEN ROUND(cost / pay_count, 4) ELSE 0 END AS cpu
		FROM @THIS@ AS t
       WHERE item != ""
	`)

	return c.RetMapping(map[string]string{
		"platform":       "渠道",
		"item":           "DD",
		"impressions":    "Impression",
		"click":          "Click",
		"ctr":            "CTR",
		"cpc":            "CPC",
		"cost":           "Spend",
		"pay_count":      "Order",
		"pay":            "Sale",
		"cart":           "Cart",
		"fav_shop_count": "Fav",
		"roi":            "ROI",
		"cvr":            "CVR",
		"cpu":            "CPU",
	})
}
