package utils

import (
	"os"
	"runtime"
	"strings"
)

//检测目录或文件是否存在, flase:不存在， true:存在
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func Testing() bool {
	wd, _ := os.Getwd()
	uintDir := "/unit"
	if runtime.GOOS == "windows" {
		uintDir = "\\unit"
	}
	return strings.Contains(wd, uintDir)
}
