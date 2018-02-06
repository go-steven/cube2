package main

import (
	"fmt"
	"github.com/go-steven/cube2/cube"
	"github.com/go-steven/cube2/engine"
	"github.com/go-steven/cube2/util"
)

func cube_join_cube() cube.Cube {
	c1 := cube.New().FromTable("skyline.simba_adgroup_rpt_daily")
	c2 := cube.New().FromTable("skyline.simba_adgroups")
	c3 := cube.New().FromTable("skyline.simba_items")

	c1.Link("@AD_GROUPS@", c2)
	c1.Link("@ITEMS@", c3)

	return c1.SQL(`SELECT
		r.*
	FROM @THIS@ AS r
   INNER JOIN @AD_GROUPS@ AS ad ON (
		ad.id = r.adgroup_id
    )
    LEFT JOIN @ITEMS@ AS item ON (
		item.id = ad.num_iid
    )
    WHERE ad.client_id = 1
    	AND r.record_on BETWEEN '2017-03-06' AND '2017-03-12'
    ORDER BY r.record_on DESC
    LIMIT 0, 10
`)
}
func main() {
	reports := engine.NewReports()
	reports.AddCube("cube join example", cube_join_cube())
	ret, err := reports.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("ret:", util.Json(ret))
}
