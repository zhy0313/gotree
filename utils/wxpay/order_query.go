package wxpay

import (
	"encoding/xml"
	"errors"
)

type OrderQuery struct {
	WeChatClient
}

// 初始化微信查询订单接口固定参数
func NewOrderQuery(host, payChannel, outTradeNo, nonceStr string) (*OrderQuery, error) {

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
	oq := &OrderQuery{}
	// 绑定配置参数
	WxHost = host
	Key = Account.Key

	oq.AppId = Account.AppId
	oq.MchId = Account.MchId
	oq.SignType = "MD5" //签名类型，默认为MD5
	oq.OutTradeNo = outTradeNo
	oq.NonceStr = nonceStr
	return oq, nil
}

// 微信查询订单接口
func (oq *OrderQuery) WeChatOrderQuery() (re OrderQueryResponse, e error) {

	oqr := OrderQueryResponse{}

	// 组合签名参数
	param := make(map[string]string, 0)
	param["appid"] = oq.AppId
	param["mch_id"] = oq.MchId
	param["out_trade_no"] = oq.OutTradeNo
	param["nonce_str"] = oq.NonceStr
	param["sign_type"] = oq.SignType

	oq.Sign = makeSign(param, Key)

	xmlParam := paramToXML(oq)

	// 执行http请求
	response, e := sendXmlRequest("POST", WxHost+"/pay/orderquery", xmlParam)
	if e != nil {
		return
	}
	// 解析xml
	e = xml.Unmarshal(response, &oqr)

	if e != nil {
		return
	} else if oqr.ReturnCode != "SUCCESS" {
		e = errors.New(oqr.ReturnMsg)
		return
	}

	re = oqr
	return
}
