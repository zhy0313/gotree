package wxpay

import (
	"testing"
)

// 微信统一下单
func TestUnifiedOrder(t *testing.T) {
	ufo, e := NewWeChatUnifiedOrder("https://api.mch.weixin.qq.com", "call", "www.baidu.com", "APP", "123445", "dfdfdfd", "1111111", "dfsfds", 0.01, "")
	if e != nil {
		t.Log(e)
	}
	t.Log(ufo.WeChatUnifiedOrder())
}

// 退款申请
func TestRefund(t *testing.T) {
	rfd, e := NewRefund("https://api.mch.weixin.qq.com", "call", "xxxxxx", "测试退款", "https://www.baidu.com", "dfsf", "dfsdfs", 0.01, 0.01)
	if e != nil {
		t.Log(e)
	}
	t.Log(rfd.WeChatRefund())
}

// 查询订单
func TestOrderQuery(t *testing.T) {
	oq, e := NewOrderQuery("https://api.mch.weixin.qq.com", "call", "1233445", "dddd")
	if e != nil {
		t.Log(e)
	}
	t.Log(oq.WeChatOrderQuery())
}

// 关闭订单
func TestOrderClose(t *testing.T) {
	oc, e := NewOrderClose("https://api.mch.weixin.qq.com", "call", "1233445", "ddd")
	if e != nil {
		t.Log(e)
	}
	t.Log(oc.WeChatOrderClose())
}

// 退款查询
func TestRefundQuery(t *testing.T) {
	rq, e := NewRefundQuery("https://api.mch.weixin.qq.com", "call", "1233445", "ddd")
	if e != nil {
		t.Log(e)
	}
	t.Log(rq.WeChatRefundQuery())
}
