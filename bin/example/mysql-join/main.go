package main

import (
	"fmt"
	"github.com/go-steven/cube2/cube"
	"github.com/go-steven/cube2/engine"
	"github.com/go-steven/cube2/util"
)

func mysql_join_cube() cube.Cube {
	return cube.New().SQL(`SELECT
		r.*
	FROM skyline.simba_adgroup_rpt_daily AS r
   INNER JOIN skyline.simba_adgroups AS ad ON (
		ad.id = r.adgroup_id
    )
    LEFT JOIN skyline.simba_items AS item ON (
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
	reports.AddCube("mysql join example", mysql_join_cube())
	ret, err := reports.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("ret:", util.Json(ret))
}
