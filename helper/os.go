package helper

import (
	"os"
	"runtime"
	"strings"
)

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
