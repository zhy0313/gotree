package wxpay

import (
	"encoding/xml"
	"errors"
)

type OrderClose struct {
	WeChatClient
}

// 初始化微信关闭订单接口固定参数
func NewOrderClose(host, payChannel, outTradeNo, nonceStr string) (*OrderClose, error) {

	var wechatAccount string
	if payChannel == "call" {
		wechatAccount = CALL_WECHAT_ACCOUNT
	} else if payChannel == "mini" {
		wechatAccount = SP_WECHAT_ACCOUNT
	}
	// 解析微信账号
	Account, e := decodeWechatAccount(wechatAccount)
	if e != nil {
		e = errors.New("解析微信账号错误")
		return nil, e
	}
	oc := &OrderClose{}
	// 绑定配置参数
	WxHost = host
	Key = Account.Key
	oc.AppId = Account.AppId
	oc.MchId = Account.MchId
	oc.SignType = "MD5" //签名类型，默认为MD5
	oc.OutTradeNo = outTradeNo
	oc.NonceStr = nonceStr
	return oc, nil
}

// 微信关闭订单接口
func (oc *OrderClose) WeChatOrderClose() (re OrderCloseResponse, e error) {

	ocr := OrderCloseResponse{}
	// 组合签名参数
	param := make(map[string]string, 0)
	param["appid"] = oc.AppId
	param["mch_id"] = oc.MchId
	param["out_trade_no"] = oc.OutTradeNo
	param["nonce_str"] = oc.NonceStr
	param["sign_type"] = oc.SignType

	oc.Sign = makeSign(param, Key)

	xmlParam := paramToXML(oc)

	// 执行http请求
	response, e := sendXmlRequest("POST", WxHost+"/pay/closeorder", xmlParam)
	if e != nil {
		return
	}
	// 解析xml
	e = xml.Unmarshal(response, &ocr)

	if e != nil {
		return
	} else if ocr.ReturnCode != "SUCCESS" {
		e = errors.New(ocr.ReturnMsg)
		return
	}

	re = ocr
	return
}
