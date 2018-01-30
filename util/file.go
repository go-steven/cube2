package util

import (
	"github.com/go-steven/cube2/util/errors"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
)

func FileExists(filename string) bool {
	var exist = true
	if fileInfo, err := os.Stat(filename); os.IsNotExist(err) || !fileInfo.Mode().IsRegular() {
		exist = false
	}
	return exist
}

func WriteFile(filename string, content []byte) error {
	return ioutil.WriteFile(filename, content, 0666)
}

func ReadFile(filename string) ([]byte, error) {
	if !FileExists(filename) {
		return nil, errors.Errorf("file[%s] not exists.", filename)
	}
	ret, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.NewErr(err)
	}

	return ret, nil
}

func RemoveFile(filename string) error {
	if !FileExists(filename) {
		return errors.Errorf("file[%s] not exists.", filename)
	}

	return os.Remove(filename)
}

func CurrDir() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "."
	}
	if runtime.GOOS == "windows" {
		file = strings.Replace(file, "\\", "/", -1)
	}
	return path.Dir(file)
}
