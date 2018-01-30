package engine

type ScriptEngine interface {
	Execute(script string, tplcfgs string) (map[string]*ReportRet, error)
}
