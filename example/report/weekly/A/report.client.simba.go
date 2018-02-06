package main

import (
	"github.com/go-steven/cube2/cube"
)

// 天猫账户整体表现（直通车）
func sima_client_period() cube.Cube {
	c := simba_client_base_period(CURR_PERIOD) // curr
	c.Link("@LAST@", simba_client_base_period(LAST_PERIOD))

	c.SQL(`SELECT record_on, impressions, click, ctr, cpc, cost, pay_count, pay, cart, fav_shop_count, roi, cvr, cpu FROM @THIS@ AS c
    UNION ALL
	SELECT record_on, impressions, click, ctr, cpc, cost, pay_count, pay, cart, fav_shop_count, roi, cvr, cpu FROM @LAST@ AS l
	`).ContrastSummary("对比", []string{"impressions", "click", "ctr", "cpc", "cost", "pay_count", "pay", "cart", "fav_shop_count", "roi", "cvr", "cpu"})

	return c.RetMapping(map[string]string{
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
