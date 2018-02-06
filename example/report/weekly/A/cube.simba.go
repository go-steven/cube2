package main

import (
	"github.com/go-steven/cube2/cube"
)

func simba_client_base_period(flag int) cube.Cube {
	periodSql := ""
	if flag == CURR_PERIOD { // curr
		periodSql = `AND record_on BETWEEN DATE_SUB('@START_DATE@', INTERVAL 1 @PERIOD@) AND DATE_SUB('@START_DATE@', INTERVAL 1 DAY)`
	} else { // last
		periodSql = `AND record_on BETWEEN '@START_DATE@' AND DATE_SUB(DATE_ADD('@START_DATE@', INTERVAL 1 @PERIOD@), INTERVAL 1 DAY)`
	}

	c := cube.New().SQL(`SELECT
			SUM(impressions) AS impressions,
			SUM(click) AS click,
			SUM(ROUND(cost/100.0,2)) AS cost,
			SUM(direct_pay_count + indirect_pay_count) AS pay_count, 
			SUM(ROUND((direct_pay + indirect_pay)/100.0,2)) AS pay,
			SUM(carttotal) AS cart,
			SUM(fav_item_count) AS fav_item_count,
			SUM(fav_shop_count) AS fav_shop_count
		FROM skyline.simba_client_rpt_daily
		WHERE client_id = @CLIENT_ID@ AND source = 0 ` + periodSql).SQL(`SELECT
			"直通车" AS platform,
			CASE '@PERIOD@' WHEN 'WEEK' THEN '上周' WHEN 'MONTH' THEN '上月' WHEN 'QUARTER' THEN '上季度' WHEN 'YEAR' THEN '上年' ELSE '' END AS record_on,
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
		"platform":       "平台",
		"record_on":      "Date",
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

func simba_product_tmp(product string) cube.Cube {
	c := cube.New().FromTable("skyline.simba_adgroup_rpt_daily")
	c.Link("@ADGROUPS@", cube.New().FromTable("skyline.simba_adgroups"))
	c.Link("@ITEMS@", cube.New().FromTable("skyline.simba_items"))
	c.Link("@CAMPAIGNS@", cube.New().FromTable("skyline.simba_campaigns"))

	return c.SQL(`SELECT 
			CASE i.title
				WHEN REPLACE(i.title, " ", "") REGEXP ".*HP01|HP02|AM11|AM10|AM06|AM07|AM08|TP02|DP01|AM09.*" THEN "ec"
				WHEN REPLACE(i.title, " ", "") REGEXP ".*DC26|DC36|DC48|DC52|V6TopDog|V6Car|V6Motorhead|V6Absolute|V8Absolute|V8Fluffy|V6Fluffy|V6Fluffy＋|V6Origin|CY22|V6Mattress|DC45.*" THEN "fc"
				WHEN REPLACE(i.title, " ", "") REGEXP ".*HD01.*" THEN "pc"
           ELSE ""
			END AS product,
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
		GROUP BY product`).SQL(`SELECT
			"直通车" AS platform,
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
       WHERE product = "` + product + `"
	`)
}
