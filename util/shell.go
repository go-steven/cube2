package util

import (
	"os/exec"
)

func ExecShell(cmd string) (string, error) {
	f, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return string(f), nil
}
