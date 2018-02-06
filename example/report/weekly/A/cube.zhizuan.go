package main

import (
	"github.com/go-steven/cube2/cube"
)

func zhizuan_client_base_period(flag int) cube.Cube {
	periodSql := ""
	if flag == CURR_PERIOD { // curr
		periodSql = `AND record_on BETWEEN DATE_SUB("@START_DATE@", INTERVAL 1 @PERIOD@) AND DATE_SUB("@START_DATE@", INTERVAL 1 DAY)`
	} else { // last
		periodSql = `AND record_on BETWEEN "@START_DATE@" AND DATE_SUB(DATE_ADD("@START_DATE@", INTERVAL 1 @PERIOD@), INTERVAL 1 DAY)`
	}
	c := cube.New()
	c.SQL(`SELECT
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
		WHERE client_id = @CLIENT_ID@ ` +  periodSql)
	c.SQL(`SELECT
			"钻展" AS platform,
			CASE '@PERIOD@' WHEN 'WEEK' THEN '上周' WHEN 'MONTH' THEN '上月' WHEN 'QUARTER' THEN '上季度' WHEN 'YEAR' THEN '上年' ELSE '' END AS record_on,
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
	`)

	return c
}

func zhizuan_product_tmp(product string) cube.Cube {
	c := cube.New().SQL(`SELECT 
			CASE campaign_name
				WHEN REPLACE(campaign_name, " ", "") REGEXP ".*EC.*" THEN "ec"
				WHEN REPLACE(campaign_name, " ", "") REGEXP ".*FC.*" THEN "fc"
				WHEN REPLACE(campaign_name, " ", "") REGEXP ".*吹风机.*" THEN "pc"
			ELSE ""
			END AS product,
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
			AND record_on BETWEEN "@START_DATE@" AND DATE_SUB(DATE_ADD("@START_DATE@", INTERVAL 1 @PERIOD@), INTERVAL 1 DAY)
		GROUP BY product`)

	return c.SQL(`SELECT
			"钻展" AS platform,
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
       WHERE product = "` + product + `"
	`)
}
