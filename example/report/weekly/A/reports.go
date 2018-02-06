package main

import (
	"github.com/go-steven/cube2/engine"
)

func Reports() *engine.Reports {
	reports := engine.NewReports()
	reports.AddCube("天猫账户整体表现（店铺）", client_curr_period())
	reports.AddCube("天猫账户整体表现（直通车）", sima_client_period())
	reports.AddCube("天猫账户整体表现（钻展）", zhizuan_client_period())
	reports.AddCube("直通车 & 钻展分类表现（EC）", cube_pruduct(EC))
	reports.AddCube("直通车 & 钻展分类表现（FC）", cube_pruduct(FC))
	reports.AddCube("直通车 & 钻展分类表现（PC）", cube_pruduct(PC))
	reports.AddCube("直通车子分类表现", simba_item())
	reports.AddCube("钻展子分类表现", zhizuan_item())

	return reports
}