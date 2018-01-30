package main

import (
	"fmt"
	"github.com/go-steven/cube2/cube"
	"github.com/go-steven/cube2/engine"
	"github.com/go-steven/cube2/util"
)

func cube_union_cube() cube.Cube {
	c1 := cube.New().SQL(`SELECT *
	FROM skyline.zhizuan_campaign_rpt_daily
    WHERE client_id = 1
    	AND record_on BETWEEN '2017-03-06' AND '2017-03-12'
    ORDER BY r.record_on DESC
    LIMIT 0, 1`)
	c2 := cube.New().SQL(`SELECT *
	FROM skyline.zhizuan_campaign_rpt_daily
    WHERE client_id = 1
    	AND record_on BETWEEN '2017-03-06' AND '2017-03-12'
    ORDER BY r.record_on DESC
    LIMIT 0, 1`)
	c1.Link("@XXX@", c2)

	return cube.New().SQL(`SELECT * FROM @THIS@
   	UNION
	SELECT * FROM @XXX@
`)
}
func main() {
	reports := engine.NewReports()
	reports.AddCube("cube union example", cube_union_cube())
	ret, err := reports.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("ret:", util.Json(ret))
}
