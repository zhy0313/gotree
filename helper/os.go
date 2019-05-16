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

func Exit(errorMsg ...string) {
	if len(errorMsg) == 0 {
		os.Exit(0)
	}
	Log().WriteError(errorMsg)
	os.Exit(-1)
}
