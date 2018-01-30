package main

import (
	"flag"
	//"fmt"
	"fmt"
	"github.com/go-steven/cube2/cube"
	"github.com/go-steven/cube2/engine"
	"github.com/go-steven/cube2/util"
	"github.com/go-steven/cube2/util/logger"
)

const (
	IMPRESSION = 1
	CLICK      = 2
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
	`).RetFieldsMapping(map[string]string{
		"record_on": "投放日期",
		"product":   "品牌",
		"cost":      "花费",
		"i_roi":     "15天展示ROI",
		"i_pay":     "15天成交金额",
	}).SummarySQL("总计", `SELECT 
		SUM(cost) AS cost, 
		SUM(i_pay) AS i_pay, 
		ROUND(IF(SUM(cost)=0, 0, SUM(i_pay) / SUM(cost)), 2) AS i_roi 
	FROM @CUBE@ AS s`).SummaryFieldsMapping("总计", map[string]string{
		"cost":  "花费",
		"i_pay": "15天成交金额",
		"i_roi": "15天展示ROI",
	})

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
		CASE WHEN cost>0 THEN ROUND(1.0 * pay / cost,2) ELSE 0 END AS roi,
		ROUND(gmv_amt,0) AS gmv_amt,
		ROUND(pay,0) AS pay,
		CASE WHEN impressions>0 THEN ROUND(100.0 * click / impressions,2) ELSE 0 END AS ctr,
		CASE WHEN click>0 THEN ROUND(1.0 * cost / click,2) ELSE 0 END AS cpc,
		CASE WHEN impressions>0 THEN ROUND(1000.0 * cost / impressions,2) ELSE 0 END AS cpm,
		CASE WHEN click>0 THEN ROUND(100.0 * pay_count / click,2) ELSE 0 END AS cvr
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
	ORDER BY i.record_on ASC`).RetFieldsMapping(map[string]string{
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
	}).SummarySQL("展示汇总", `SELECT
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
	`).SummaryFieldsMapping("展示汇总", map[string]string{
		"cost":           "全店花费",
		"impressions":    "Impression",
		"click":          "Click",
		"pay_count":      "Orders",
		"fav_shop_count": "Store Collection",
		"fav_item_count": "Item Collection",
		"cart":           "Add Cart 15 Days",
		"gmv_amt":        "15订单金额",
		"i_pay":          "展现成交金额",
		"c_pay":          "点击成交金额",
		"ctr":            "CTR",
		"cpc":            "CPC",
		"cpm":            "CPM",
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
		CASE WHEN cost>0 THEN ROUND(1.0 * pay / cost,2) ELSE 0 END AS roi,
		ROUND(gmv_amt,0) AS gmv_amt,
		ROUND(pay,0) AS pay,
		CASE WHEN impressions>0 THEN ROUND(100.0 * click / impressions,2) ELSE 0 END AS ctr,
		CASE WHEN click>0 THEN ROUND(1.0 * cost / click,2) ELSE 0 END AS cpc,
		CASE WHEN impressions>0 THEN ROUND(1000.0 * cost / impressions,2) ELSE 0 END AS cpm,
		CASE WHEN click>0 THEN ROUND(100.0 * pay_count / click,2) ELSE 0 END AS cvr
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

var (
	Logger        = logger.Logger
	tplConfigFlag = flag.String("tplcfg", "", "set tpl config file, .json")
	outputFlag    = flag.String("output", "", "set reports json output")
)

func main() {
	flag.Parse()

	reports := engine.NewReports()
	clientIds := []uint64{10, 9, 8, 6, 7, 17, 5}
	// 所有客户汇总统计报表
	reports.AddCube("client_zhizuan_daily_report_0_0", all_clients_stats(clientIds))
	// 单个客户统计报表
	for k, v := range clientIds {
		reports.AddCube(fmt.Sprintf("client_zhizuan_daily_report_%d_c%d", k+1, v), client_stats(v))
	}

	if err := reports.RunAndSave(*tplConfigFlag, *outputFlag); err != nil {
		Logger.Error(err)
		return
	}
}
