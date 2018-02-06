package engine

import (
	"github.com/go-steven/cube2/cube"
	"github.com/go-steven/cube2/source"
	"github.com/go-steven/cube2/util"
	"github.com/go-steven/cube2/util/errors"
	"sync"
)

type ReportRet struct {
	Name    string                `json:"-", yaml:"-"`
	Display interface{}           `json:"display", yaml:"display,omitempty"` // 在结果中原样返回，用户前段显示数据用
	Fields  []string              `json:"fields", yaml:"fields,omitempty"`
	Data    source.Rows           `json:"data", yaml:"data,omitempty"`
	Summary map[string]source.Row `json:"summary", yaml:"data,omitempty"`
}

type Reports struct {
	data map[string]cube.Cube
	m    *sync.RWMutex // 用于并发
}

func NewReports() *Reports {
	return &Reports{
		data: make(map[string]cube.Cube),
		m:    new(sync.RWMutex),
	}
}

func (r *Reports) AddCube(name string, c cube.Cube) *Reports {
	name = util.Trim(name)
	if name != "" && c != nil {
		r.m.Lock()
		r.data[name] = c
		r.m.Unlock()
	}
	return r
}

func (r *Reports) Cubes() map[string]cube.Cube {
	ret := make(map[string]cube.Cube)
	r.m.RLock()
	for k, v := range r.data {
		if v == nil {
			continue
		}
		ret[k] = v
	}
	r.m.RUnlock()
	return ret
}

func (r *Reports) Run() (map[string]*ReportRet, error) {
	tplCfgs := make(cube.TplCfg)
	return r.RunWithCfgs(tplCfgs)
}

// 生成各个报表的数据，并发执行
func (r *Reports) RunWithCfgs(tplCfgs cube.TplCfg) (map[string]*ReportRet, error) {
	ret := make(map[string]*ReportRet)
	m := new(sync.Mutex)  // 用于并发
	var wg sync.WaitGroup // 用于并发
	for k, v := range r.Cubes() {
		if v == nil {
			continue
		}

		wg.Add(1)
		go func(name string, c cube.Cube) error {
			defer wg.Done()

			report := &ReportRet{
				Name: name,
			}
			c.Replace(tplCfgs)

			fields, err := c.Fields()
			if err != nil {
				Logger.Error(err)
				return err
			}
			rows, err := c.Rows()
			if err != nil {
				Logger.Error(err)
				return err
			}
			report.Data = rows
			report.Fields = fields

			summary, err := c.Summary()
			if err != nil {
				Logger.Error(err)
				return err
			}
			if summary != nil {
				report.Summary = summary
			}
			m.Lock()
			ret[name] = report
			m.Unlock()

			return nil
		}(k, v.Copy())
	}
	wg.Wait()
	Logger.Infof("Run report: %s", util.Json(ret))
	return ret, nil
}

func (r *Reports) RunAndSave(tplConfigFile, outputFile string) error {
	if outputFile == "" {
		return errors.New("No output file.")
	}
	tplcfg := make(cube.TplCfg)
	if tplConfigFile != "" {
		var err error
		tplcfg, err = cube.ReadTplCfgFile(tplConfigFile)
		if err != nil {
			return errors.NewErr(err)
		}
	}
	reports, err := r.RunWithCfgs(tplcfg)
	if err != nil {
		return errors.NewErr(err)
	}
	if err := util.WriteFile(outputFile, []byte(util.Json(reports))); err != nil {
		return errors.NewErr(err)
	}

	return nil
}
