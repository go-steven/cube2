package main

import (
	"github.com/go-steven/cube2/cube"
)

// 直通车 & 钻展分类表现（EC/FC/PC）
func cube_pruduct(product string) cube.Cube {
	c := simba_product_tmp(product)
	c.Link("@ZHIZUAN@", zhizuan_product_tmp(product))

	c.SQL(`SELECT platform, impressions, click, ctr, cpc, cost, pay_count, cart, fav_shop_count, roi, cvr, cpu FROM @THIS@ AS s
    UNION ALL
	SELECT platform, impressions, click, ctr, cpc, cost, pay_count, cart, fav_shop_count, roi, cvr, cpu FROM @ZHIZUAN@ AS z
	`)

	return c.RetMapping(map[string]string{
		"platform":       "平台",
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
