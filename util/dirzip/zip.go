package dirzip

import (
	"encoding/json"
	"github.com/go-steven/cube2/util"
	"github.com/go-steven/cube2/util/errors"
	"io/ioutil"
	"path"
	"strings"
)

func Zip(dirname string, addMain bool) (string, error) {
	ok, err := util.PathExists(dirname)
	if err != nil {
		return "", errors.NewErr(err)
	}
	if !ok {
		return "", errors.Errorf("dir not exists: %s", dirname)
	}
	dirname = path.Clean(dirname) + "/"

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return "", err
	}

	content_map := make(map[string]string)
	for _, v := range files {
		if !v.Mode().IsRegular() {
			continue // IGNORE
		}

		ext := path.Ext(v.Name())
		if ext != ".go" {
			continue // IGNORE
		}

		content, err := util.ReadFile(dirname + v.Name())
		if err != nil {
			return "", err
		}
		content_map[v.Name()] = string(content)
	}

	if addMain {
		if v, ok := content_map["test.go"]; ok {
			content_map["test.go"] = strings.Replace(v, "func main()", "func main_test()", -1)
		}
		content, err := readMainExample()
		if err != nil {
			return "", err
		}
		content_map["main.go"] = content
	}
	return util.Json(content_map), nil
}

func UnZip(zip string, dirname string, addMain bool) error {
	ok, err := util.PathExists(dirname)
	if err != nil {
		return errors.NewErr(err)
	}
	if !ok {
		return errors.Errorf("dir not exists: %s", dirname)
	}
	dirname = path.Clean(dirname) + "/"

	content_map := make(map[string]string)
	if err := json.Unmarshal([]byte(zip), &content_map); err != nil {
		return errors.Errorf("ERROR Unmarshal: %v", err.Error())
	}

	if addMain {
		content, err := readMainExample()
		if err != nil {
			return err
		}
		content_map["main.go"] = content
	} else {
		if v, ok := content_map["test.go"]; ok {
			content_map["test.go"] = strings.Replace(v, "func main_test()", "func main()", -1)
		}
	}

	for k, v := range content_map {
		if err := util.WriteFile(dirname+k, []byte(v)); err != nil {
			return errors.NewErr(err)
		}
	}

	return nil
}

func readMainExample() (string, error) {
	content, err := util.ReadFile(util.CurrDir() + "/" + "main.go.example")
	if err != nil {
		return "", errors.NewErr(err)
	}
	return string(content), nil
}
