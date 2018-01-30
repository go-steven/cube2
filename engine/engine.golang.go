package engine

import (
	"encoding/json"
	"fmt"
	"github.com/go-steven/cube2/util"
	"github.com/go-steven/cube2/util/errors"
	"os"
)

const (
	DEFAULT_TMP_DIR = "/tmp/"
)

type GoEngine struct {
	tmpdir string
}

func NewGoEngine() *GoEngine {
	return &GoEngine{
		tmpdir: getDefaultTmpDir(),
	}
}

func (e *GoEngine) Execute(script string, tplcfgs string) (map[string]*ReportRet, error) {
	Logger.Infof("Execute:")
	Logger.Info(script)
	Logger.Infof("tplcfgs: %s", tplcfgs)
	// save script to .go file
	scriptFile := fmt.Sprintf("%s%s_main.go", e.tmpdir, util.Token())
	Logger.Infof("script file: %s", scriptFile)
	if err := util.WriteFile(scriptFile, []byte(script)); err != nil {
		return nil, errors.NewErr(err)
	}
	// save tplcfgs to .cfg file
	cfgFile := fmt.Sprintf("%s%s_tpl.cfg", e.tmpdir, util.Token())
	Logger.Infof("cfg file: %s", cfgFile)
	if err := util.WriteFile(cfgFile, []byte(tplcfgs)); err != nil {
		return nil, errors.NewErr(err)
	}

	outputFile := fmt.Sprintf("%s%s.output", e.tmpdir, util.Token())
	Logger.Infof("output file: %s", outputFile)
	// run shell and save the result to a file
	cmd := fmt.Sprintf(`go run %s --tplcfg=%s --output=%s`, scriptFile, cfgFile, outputFile)
	Logger.Infof("run shell: %s", cmd)
	_, err := util.ExecShell(cmd)
	if err != nil {
		return nil, errors.NewErr(err)
	}
	// read the result file and parse it
	output, err := util.ReadFile(outputFile)
	if err != nil {
		return nil, errors.NewErr(err)
	}
	//Logger.Infof("read output: %s", string(output))
	reports := make(map[string]*ReportRet)
	if err := json.Unmarshal(output, &reports); err != nil {
		Logger.Errorf("ERROR Unmarshal: %v", errors.NewErr(err))
		return nil, err
	}
	//Logger.Infof("reports: %s", util.Json(reports))
	// cleanup
	if err := util.RemoveFile(scriptFile); err != nil {
		Logger.Error(errors.NewErr(err))
	}
	if err := util.RemoveFile(cfgFile); err != nil {
		Logger.Error(errors.NewErr(err))
	}
	if err := util.RemoveFile(outputFile); err != nil {
		Logger.Error(errors.NewErr(err))
	}

	return reports, nil
}

func getDefaultTmpDir() string {
	fileInfo, err := os.Stat(DEFAULT_TMP_DIR)
	if err == nil && fileInfo.IsDir() {
		return DEFAULT_TMP_DIR
	}

	return "./"
}
