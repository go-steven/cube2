package util

import (
	"crypto/sha1"
	"fmt"
	"io"
	"time"
)

func Sha1(s string) string {
	t := sha1.New()
	io.WriteString(t, s)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func Sha1Name(name string) string {
	return fmt.Sprintf("t_%s", Sha1(name))
}

func Token() string {
	return Sha1(time.Now().Format("2006-01-02 15:04:05.999999999 -0700 MST"))
}
