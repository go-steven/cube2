package main

import (
	"github.com/go-steven/cube2/cube"
)

// 直通车子分类表现
func simba_item() cube.Cube {
	c := cube.New().FromTable("skyline.simba_adgroup_rpt_daily")
	c.Link("@ADGROUPS@", cube.New().FromTable("skyline.simba_adgroups"))
	c.Link("@ITEMS@", cube.New().FromTable("skyline.simba_items"))
	c.Link("@CAMPAIGNS@", cube.New().FromTable("skyline.simba_campaigns"))

	c.SQL(`SELECT 
			CASE i.title
				WHEN REPLACE(i.title, " ", "") REGEXP ".*AM11|HP01|HP02.*" THEN "净化"
				WHEN REPLACE(i.title, " ", "") REGEXP ".*AM06|AM07|AM08|TP02|DP01.*" THEN "风扇"
				WHEN REPLACE(i.title, " ", "") REGEXP ".*AM09.*" THEN "取暖器"
				WHEN REPLACE(i.title, " ", "") REGEXP ".*AM10.*" THEN "加湿器"
				WHEN REPLACE(i.title, " ", "") REGEXP ".*V6TopDog|V6Car|V6Motorhead|V6Origin|V6Fluffy|V6Fluffy＋|V6Absolute|V8Absolute|V8Fluffy|DC45.*" 
					 AND REPLACE(i.title, " ", "") NOT REGEXP ".*车载吸尘器|圆筒吸尘器.*" 
				THEN "吸尘器"
				WHEN REPLACE(i.title, " ", "") REGEXP ".*V6Car.*" THEN "车载吸尘器"
				WHEN REPLACE(i.title, " ", "") REGEXP ".*CY22|DC26|DC36|DC48|DC52.*" THEN "圆筒吸尘器"
				WHEN REPLACE(i.title, " ", "") REGEXP ".*V6Mattress.*" THEN "除螨仪"
				WHEN REPLACE(i.title, " ", "") REGEXP ".*HD01.*" THEN "吹风机"
           ELSE ""
			END AS item,
			SUM(impressions) AS impressions,
			SUM(click) AS click,
			SUM(ROUND(cost/100.0,2)) AS cost,
			SUM(direct_pay_count + indirect_pay_count) AS pay_count, 
			SUM(ROUND((direct_pay + indirect_pay)/100.0,2)) AS pay,
			SUM(carttotal) AS cart,
			SUM(fav_item_count) AS fav_item_count,
			SUM(fav_shop_count) AS fav_shop_count
		FROM @THIS@ AS r
		INNER JOIN @ADGROUPS@ AS ad ON ad.id = r.adgroup_id
		INNER JOIN @ITEMS@ AS i ON i.id = ad.num_iid
		INNER JOIN @CAMPAIGNS@ AS c ON c.id = ad.campaign_id
		WHERE r.client_id = @CLIENT_ID@ AND r.source = 0
			AND r.record_on BETWEEN '@START_DATE@' AND DATE_SUB(DATE_ADD('@START_DATE@', INTERVAL 1 @PERIOD@), INTERVAL 1 DAY)
           AND c.title NOT LIKE '%皇冠车手%'
			AND c.title NOT LIKE '%阻击计划%'
			AND c.title NOT LIKE '%干手器%'
			AND c.title NOT LIKE '%配件%'
		GROUP BY record_on`).SQL(`SELECT
			"直通车" AS platform,
			item,
			impressions,
			click,
			CASE WHEN impressions>0 THEN ROUND(click / impressions, 4) ELSE 0 END AS ctr,
			CASE WHEN click>0 THEN ROUND(cost / click, 4) ELSE 0 END AS cpc,
			ROUND(cost,2) AS cost,
			pay_count,
			ROUND(pay,2) AS pay,
			cart,
			fav_shop_count,
			fav_item_count,
			CASE WHEN cost>0 THEN ROUND(pay / cost, 4) ELSE 0 END AS roi,
			CASE WHEN click>0 THEN ROUND(pay_count / click, 4) ELSE 0 END AS cvr,
			CASE WHEN pay_count>0 THEN ROUND(cost / pay_count, 4) ELSE 0 END AS cpu
		FROM @THIS@ AS t
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
