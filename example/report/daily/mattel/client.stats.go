package main

import (
	"github.com/go-steven/cube2/cube"
)

/***************************************************************************************************************
***
***   单个客户统计报表相关CUBE
***
***************************************************************************************************************
 */
// 客户推广统计报表：全店推广 + 单品推广
// effectType: 效果类型。1:“impression”：展现效果；2:“click”：点击效果
func client_stats(clientId uint64) cube.Cube {
	// 展现效果统计报表
	impression_stats := client_promotion_stats(clientId, IMPRESSION)
	// 关联点击效果统计报表
	impression_stats.Link("@CLICK@", client_promotion_stats(clientId, CLICK))

	return impression_stats.SQL(`SELECT 
		i.record_on,
		skyline.weekinfo(i.record_on) AS week_info, 
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
	LEFT JOIN @CLICK@ AS c ON c.record_on = i.record_on 
	ORDER BY i.record_on ASC`).SummarySQL("展示汇总", `SELECT
		SUM(cost) AS cost,
		SUM(impressions) AS impressions,
		SUM(click) AS click,
		SUM(pay_count) AS pay_count,
		SUM(fav_shop_count) AS fav_shop_count,
		SUM(fav_item_count) AS fav_item_count,
		SUM(cart) AS cart,
		SUM(gmv_amt) AS gmv_amt,
		SUM(i_pay) AS i_pay,
		SUM(c_pay) AS c_pay,
		ROUND(IF(SUM(click)=0, 0, SUM(cost)/SUM(click)),2) AS cpc,
		ROUND(100*IF(SUM(impressions)=0, 0, SUM(click)/SUM(impressions)),2) AS ctr,
		ROUND(1000*IF(SUM(impressions)=0, 0, SUM(cost)/SUM(impressions)),2) AS cpm
	FROM @CUBE@ AS r
	`).RetMapping(map[string]string{
		"record_on":      "Date",
		"week_info":      "周数据",
		"cost":           "全店花费",
		"impressions":    "Impression",
		"click":          "Click",
		"pay_count":      "Orders",
		"fav_shop_count": "Store Collection",
		"fav_item_count": "Item Collection",
		"cart":           "Add Cart 15 Days",
		"i_roi":          "全店展示ROI",
		"c_roi":          "全店点击ROI",
		"gmv_amt":        "15订单金额",
		"i_pay":          "展现成交金额",
		"c_pay":          "点击成交金额",
		"ctr":            "CTR",
		"cpc":            "CPC",
		"cpm":            "CPM",
		"cvr":            "CVR",
	})
}

// 客户推广统计报表：全店推广 + 单品推广
// effectType: 效果类型。1:“impression”：展现效果；2:“click”：点击效果
func client_promotion_stats(clientId uint64, effectType uint8) cube.Cube {
	// 全店推广统计
	shop_promotion_stats := client_shop_promotion_stats(clientId, effectType)
	// 关联单品推广统计
	shop_promotion_stats.Link("@PRODUCT@", client_product_promotion_stats(clientId, effectType))

	return shop_promotion_stats.SQL(`SELECT  
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
	FROM (
		SELECT 
			record_on, 
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
			record_on, 
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
    GROUP BY record_on
`).SQL(`SELECT 
		record_on,
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

// 根据智钻推广计划日统计数据，得到客户全店推广日统计报表
// effectType: 效果类型。1:“impression”：展现效果；2:“click”：点击效果
func client_shop_promotion_stats(clientId uint64, effectType uint8) cube.Cube {
	return cube.New().SQL(`	SELECT 
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
		FROM skyline.zhizuan_campaign_rpt_daily_new 
		WHERE client_id = %d 
			AND campaign_name NOT LIKE 'tr%%' 
			AND record_on BETWEEN '@START_DATE@' AND '@END_DATE@' 
			AND effect = 15 
			AND effect_type = %d 
			AND campaign_model = 1 
		GROUP BY record_on`, clientId, effectType)
}

// 根据智钻店铺日统计数据，得到客户单品推广日统计报表
// effectType: 效果类型。1:“impression”：展现效果；2:“click”：点击效果
func client_product_promotion_stats(clientId uint64, effectType uint8) cube.Cube {
	return cube.New().SQL(`	SELECT 
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
		FROM skyline.zhizuan_client_rpt_daily_new 
		WHERE client_id = %d 
			AND record_on BETWEEN '@START_DATE@' AND '@END_DATE@' 
			AND effect = 15 
			AND effect_type = %d 
			AND campaign_model = 4 
		GROUP BY record_on`, clientId, effectType)
}
