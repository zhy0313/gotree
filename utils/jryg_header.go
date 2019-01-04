package utils

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

//HeaderContent json
type HeaderContent struct {
	LocLat         string `json:"loc_lat"`   // 纬度
	LocLng         string `json:"loc_lng"`   //经度
	LocType        string `json:"loc_type"`  //经纬度类型：WGS84-􏰦􏰧􏰨􏰩􏰪 GCJ02-􏰫􏰬􏰨􏰩􏰪BD09-􏰭􏰥􏰨􏰩 􏰮
	LocSpeed       string `json:"loc_speed"` //速度
	ImeiUUID       string `json:"imei_uuid"` //uuid
	NetType        string `json:"net_type"`  //2G/3G/4G/wifi
	DeviceNo       string `json:"device_no"` //􏰳􏰴Android 􏰵􏰶 key:device_no value:jkdfj99348938493
	UserID         string `json:"user_id"`
	CityID         string `json:"city_id"`
	TimeStampUnix  string `json:"timestamp_unix"`
	AppTimeoutMs   string `json:"app_timeout_ms"` //􏰻􏰷􏰌􏰼 key:app_timeout_ms value:30000
	Imsi           string `json:"imsi"`
	Mac            string `json:"mac"`
	MobileMerchant string `json:"mobileMerchant"`
}

//UserAgent 用户代理
type UserAgent struct {
	Client        string
	AppVersion    string
	DeviceType    string
	SystemVersion string
	DeviceName    string
	Channel       string
	AgentID       int //旧渠道编号
}

//Unmarshal 解析 header-content
func (hc *HeaderContent) Unmarshal(base64AppHeader string) *HeaderContent {
	decodedJSONByte, _ := base64.StdEncoding.DecodeString(base64AppHeader)
	appHeader := string(decodedJSONByte)
	json.Unmarshal([]byte(appHeader), hc)
	return hc
}

//DecodeUA 解析 User-Agent
func (ua *UserAgent) DecodeUA(userAgent string) *UserAgent {
	arr := strings.Split(userAgent, "/")
	if len(arr) > 1 {
		agentMap := map[string]bool{"YGPassenger": true, "YGDriver": true, "YGDTaxi": true, "YGGuider": true, "YGSmallProgram": true}
		if !agentMap[arr[0]] { //不是正确的请求来源
			return ua
		}
		ua.Client = arr[0] // e.g. YGPassenger
		tmp := strings.Trim(arr[1], ")")
		arr2 := strings.Split(tmp, "(")
		if len(arr2) > 1 {
			ua.AppVersion = arr2[0] // e.g. App版本号 3.1.0/3.4.0
			arr3 := strings.Split(arr2[1], ";")
			if len(arr3) > 3 {
				ua.DeviceType = strings.ToLower(arr3[0]) // e.g. 终端 Android/iOS/WebApp/WechatApp
				ua.SystemVersion = arr3[1]               // e.g. 系统版本号 android7.0/ios11.3.1/chrome12.1/Safari9.1
				ua.DeviceName = arr3[2]                  // e.g. 设备名称 huaweiP9/oppoR15/IphoneX
				ua.Channel = arr3[3]                     // 渠道号
			}
		}
		if ua.DeviceType == "android" {
			ua.AgentID = 12
		} else if ua.DeviceType == "ios" {
			ua.AgentID = 11
		} else if ua.DeviceType == "webapp" {
			ua.AgentID = 14
		} else if ua.DeviceType == "wechatapp" {
			ua.AgentID = 149
		}
	}
	return ua
}
