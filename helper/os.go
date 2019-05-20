package helper

import (
	"errors"
	"os"
	"strings"
)

var ErrBreaker = errors.New("Fuse")

type VoidValue struct {
	Void byte
}

func Testing() bool {
	wd, _ := os.Getwd()
	if strings.Contains(wd, "/unit") || strings.Contains(wd, "\\unit") {
		return true
	}
	return false
}

func Exit(errorMsg ...string) {
	if len(errorMsg) == 0 {
		os.Exit(0)
	}
	Log().Error(errorMsg)
	os.Exit(-1)
}
