package helper

import (
	"errors"
	"os"
	"runtime"
	"strings"
)

var ErrBreaker = errors.New("熔断")

type VoidValue struct {
	Void byte
}

func Testing() bool {
	wd, _ := os.Getwd()
	uintDir := "/unit"
	if runtime.GOOS == "windows" {
		uintDir = "\\unit"
	}
	return strings.Contains(wd, uintDir)
}
