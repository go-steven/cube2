package main

import (
	"github.com/go-steven/cube2/cube"
	"github.com/go-steven/cube2/util"
)

/***************************************************************************************************************
***
***   所有客户汇总统计报表相关CUBE
***
***************************************************************************************************************
 */
// 所有客户汇总推广统计报表：全店推广 + 单品推广
// effectType: 效果类型。1:“impression”：展现效果；2:“click”：点击效果
func all_clients_stats(clientIds []uint64) cube.Cube {
	if len(clientIds) == 0 {
		clientIds = append(clientIds, 0)
	}

	// 展现效果统计报表
	impression_stats := all_clients_promotion_stats(clientIds, IMPRESSION)
	// 关联点击效果统计报表
	impression_stats.Link("@CLICK@", all_clients_promotion_stats(clientIds, CLICK))

	impression_stats.SQL(`SELECT 
		CONCAT(date_format('@START_DATE@', '%m.%d'), '-', date_format('@END_DATE@', '%m.%d')) AS record_on,
		CASE i.client_id 
			WHEN 10 THEN 1 
			WHEN 9 THEN 2 
			WHEN 8 THEN 3 
			WHEN 6 THEN 4 
			WHEN 7 THEN 5 
			WHEN 17 THEN 6 
			WHEN 5 THEN 7 
		END AS order_seq,
		CASE i.client_id 
			WHEN 10 THEN 'Fisherprice' 
			WHEN 9 THEN 'Thomas&Friends' 
			WHEN 8 THEN 'Barbie' 
			WHEN 6 THEN 'Hotwheels' 
			WHEN 7 THEN 'Megabloks' 
			WHEN 17 THEN 'Octonauts' 
			WHEN 5 THEN 'Mattel' 
		END AS product,
		i.cost,
		i.impressions,
		i.click, 
		i.pay_count,
		i.fav_shop_count,
		i.fav_item_count,
		i.cart,
		i.roi AS i_roi,
		c.roi AS c_roi, 
		i.gmv_amt,
		i.pay AS i_pay,
		c.pay AS c_pay, 
		i.ctr,
		i.cpc,
		i.cpm,
		i.cvr
	FROM @THIS@ AS i
	LEFT JOIN @CLICK@ AS c ON c.client_id = i.client_id 
	ORDER BY i.client_id ASC`).SQL(`SELECT
		record_on,
		product,
		cost,
		i_roi,
		i_pay
		FROM @THIS@ AS r
		ORDER BY order_seq ASC
	`).RetMapping(map[string]string{
		"record_on": "投放日期",
		"product":   "品牌",
		"cost":      "花费",
		"i_roi":     "15天展示ROI",
		"i_pay":     "15天成交金额",
	}).SummarySQL("总计", `SELECT 
		SUM(cost) AS cost, 
		SUM(i_pay) AS i_pay, 
		ROUND(IF(SUM(cost)=0, 0, SUM(i_pay) / SUM(cost)), 2) AS i_roi 
	FROM @CUBE@ AS s`)

	return impression_stats
}

// 所有客户汇总推广统计报表：全店推广 + 单品推广
// effectType: 效果类型。1:“impression”：展现效果；2:“click”：点击效果
func all_clients_promotion_stats(clientIds []uint64, effectType uint8) cube.Cube {
	if len(clientIds) == 0 {
		clientIds = append(clientIds, 0)
	}

	// 全店推广统计
	shop_promotion_stats := all_clients_shop_promotion_stats(clientIds, effectType)
	// 关联单品推广统计
	shop_promotion_stats.Link("@PRODUCT@", all_clients_product_promotion_stats(clientIds, effectType))

	return shop_promotion_stats.SQL(`SELECT  
		client_id, 
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
	FROM (
		SELECT 
			client_id, 
			impressions, 
			uv, 
			click, 
			cost, 
			cart, 
			fav_item_count, 
			fav_shop_count,
			gmv_amt,
			gmv_count,
			pay,
			pay_count
		FROM @THIS@ AS t 
		UNION ALL 
		SELECT 
			client_id, 
			impressions, 
			uv, 
			click, 
			cost, 
			cart, 
			fav_item_count, 
			fav_shop_count,
			gmv_amt,
			gmv_count,
			pay,
			pay_count
		FROM @PRODUCT@ AS t
	) AS s
    GROUP BY client_id
`).SQL(`SELECT 
		client_id,
		ROUND(cost,0) AS cost,
		impressions,
		click,
		pay_count,
		fav_shop_count,
		fav_item_count,
		cart,
		IF(cost=0, 0, ROUND(1.0 * pay / cost,2)) AS roi,
		ROUND(gmv_amt,0) AS gmv_amt,
		ROUND(pay,0) AS pay,
		IF(impressions=0, 0, ROUND(100.0 * click / impressions,2)) AS ctr,
		IF(click=0, 0, ROUND(1.0 * cost / click,2)) AS cpc,
		IF(impressions=0, 0, ROUND(1000.0 * cost / impressions,2)) AS cpm,
		IF(click=0, 0, ROUND(100.0 * pay_count / click,2)) AS cvr
    FROM @THIS@ AS t`)
}

// 根据智钻推广计划日统计数据，得到所有客户汇总全店推广日统计报表
// effectType: 效果类型。1:“impression”：展现效果；2:“click”：点击效果
func all_clients_shop_promotion_stats(clientIds []uint64, effectType uint8) cube.Cube {
	if len(clientIds) == 0 {
		clientIds = append(clientIds, 0)
	}
	c := cube.New()
	return c.SQL(`SELECT 
			client_id,
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
		FROM skyline.zhizuan_campaign_rpt_daily_new 
		WHERE client_id IN (%s) 
			AND campaign_name NOT LIKE 'tr%%' 
			AND record_on BETWEEN '@START_DATE@' AND '@END_DATE@' 
			AND effect = 15 
			AND effect_type = %d 
			AND campaign_model = 1 
		GROUP BY client_id`, c.Escape(util.Uint64Join(clientIds)), effectType)
}

// 根据智钻店铺日统计数据，得到所有客户汇总单品推广日统计报表
// effectType: 效果类型。1:“impression”：展现效果；2:“click”：点击效果
func all_clients_product_promotion_stats(clientIds []uint64, effectType uint8) cube.Cube {
	if len(clientIds) == 0 {
		clientIds = append(clientIds, 0)
	}
	c := cube.New()
	return c.SQL(`SELECT 
			client_id,
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
		FROM skyline.zhizuan_client_rpt_daily_new 
		WHERE client_id IN (%s)
			AND record_on BETWEEN '@START_DATE@' AND '@END_DATE@' 
			AND effect = 15 
			AND effect_type = %d 
			AND campaign_model = 4 
		GROUP BY client_id`, c.Escape(util.Uint64Join(clientIds)), effectType)
}
