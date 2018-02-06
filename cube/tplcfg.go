package cube

import (
	"encoding/json"
	"github.com/go-steven/cube2/util"
	"github.com/go-steven/cube2/util/errors"
)

type TplCfg map[string]interface{}

func ReadTplCfgFile(tplCfgFile string) (TplCfg, error) {
	content, err := util.ReadFile(tplCfgFile)
	if err != nil {
		return nil, errors.NewErr(err)
	}
	tplCfgs := make(TplCfg)
	if util.Trim(string(content)) == "" {
		return tplCfgs, nil
	}
	if err := json.Unmarshal(content, &tplCfgs); err != nil {
		return nil, errors.Errorf("ERROR Unmarshal: %v", err.Error())
	}
	return tplCfgs, nil
}
