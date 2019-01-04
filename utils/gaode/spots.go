package gaode

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"jryghq.cn/utils"
	"sort"
	"strconv"
	"strings"
	"time"
)

//Spots sdk
type Spots struct {
	Key    string
	Center string
	Radius int
	Count  int
	Ts     string
	Scode  string
}

//Spot 初始化
func (s *Spots) Spot(center string, radius, count int) *Spots {
	s.Key = "79d13a44d4a4c738e1403de1d3e6a4b7"
	s.Radius = radius
	s.Count = count
	s.Center = center
	return s
}

//GetTs 计算 Ts
func (s *Spots) GetTs() *Spots {
	timestamp := time.Now().Unix()
	timestampStr := strconv.FormatInt(timestamp, 10)
	t := len(timestampStr)
	s.Ts = s.substr(timestampStr, 0, t-2) + "0" + s.substr(timestampStr, t-2, t-1)
	return s
}

//GetScode 计算 Scode
func (s *Spots) GetScode() *Spots {
	var params []string
	params = append(params, "F6:84:4D:82:FE:40:DB:74:F4:E9:F0:A7:FC:7A:97:A7:20:96:7A:C3:com.jryg.client")
	params = append(params, s.Ts)
	params = append(params, s.SortKeys())

	sign := md5.New()
	sign.Write([]byte(strings.Join(params, ":")))
	s.Scode = hex.EncodeToString(sign.Sum(nil))
	return s
}

//SortKeys 字典序
func (s *Spots) SortKeys() string {
	m := make(map[string]string, 0)

	m["center"] = s.Center
	m["radius"] = fmt.Sprintf("%d", s.Radius)
	m["count"] = fmt.Sprintf("%d", s.Count)
	m["key"] = s.Key

	sortedKeys := make([]string, 0)
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys) //字典排序
	var paramsString []string
	for _, k := range sortedKeys {
		v := m[k]
		if v != "" {
			paramsString = append(paramsString, fmt.Sprintf("%s=%s", k, v))
		}
	}
	paramsStr := strings.Join(paramsString, "&") //对参数按照key=value的格式，并按照参数名ASCII字典序排序
	return paramsStr
}

//Substr 字符串截取
func (s *Spots) substr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)
	if start < 0 || start > length {
		utils.Log().WriteError("start is wrong")
	}
	if end < 0 || end > length {
		utils.Log().WriteError("end is wrong")
	}
	return string(rs[start:end])
}
