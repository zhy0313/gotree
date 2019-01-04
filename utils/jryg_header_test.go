package utils

import (
	"testing"
)

func TestUserAgent_DecodeUA(t *testing.T) {
	ua := new(UserAgent)
	uaStr := "YGPassenger/3.9.0(WebApp;iOS12.1.3;iPhone9,1;appstore)"
	uaStrFake := "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Mobile Safari/537.36"
	ua.DecodeUA(uaStr)
	t.Log(ua)

	ua1 := new(UserAgent)
	ua1.DecodeUA(uaStrFake)
	t.Log(ua1)

}
