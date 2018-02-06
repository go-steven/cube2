package main

import (
	"fmt"
	"github.com/go-steven/cube2/engine"
)

/*
美泰日报汇总
*/
func Reports() *engine.Reports {
	reports := engine.NewReports()
	clientIds := []uint64{10, 9, 8, 6, 7, 17, 5}
	// 所有客户汇总统计报表
	reports.AddCube("client_zhizuan_daily_report_0_0", all_clients_stats(clientIds))
	// 单个客户统计报表
	for k, v := range clientIds {
		reports.AddCube(fmt.Sprintf("client_zhizuan_daily_report_%d_c%d", k+1, v), client_stats(v))
	}

	return reports
}
