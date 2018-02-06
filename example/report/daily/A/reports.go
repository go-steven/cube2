/*
通用模板-日报A
*/
package main

import (
	"github.com/go-steven/cube2/engine"
)

func Reports() *engine.Reports {
	r := engine.NewReports()
	r.AddCube("client_simba_daily_report", simba_daily_report())
	r.AddCube("client_zhizuan_daily_report", zhizuan_daily_report())

	return r
}
