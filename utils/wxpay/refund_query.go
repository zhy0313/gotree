package wxpay

import (
	"encoding/xml"
	"errors"
)

type RefundQuery struct {
	WeChatClient
}

// 初始化微信关闭订单接口固定参数
func NewRefundQuery(host, payChannel, outTradeNo, nonceStr string) (*RefundQuery, error) {

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
	rq := &RefundQuery{}
	// 绑定配置参数
	WxHost = host
	Key = Account.Key
	rq.AppId = Account.AppId
	rq.MchId = Account.MchId
	rq.SignType = "MD5" //签名类型，默认为MD5
	rq.OutTradeNo = outTradeNo
	rq.NonceStr = nonceStr
	return rq, nil
}

// 微信退款查询接口
func (rfq *RefundQuery) WeChatRefundQuery() (re RefundQueryResponse, e error) {

	rqr := RefundQueryResponse{}

	// 组合请求参数

	// 组合签名参数
	param := make(map[string]string, 0)
	param["appid"] = rfq.AppId
	param["mch_id"] = rfq.MchId
	param["out_trade_no"] = rfq.OutTradeNo
	param["nonce_str"] = rfq.NonceStr
	param["sign_type"] = rfq.SignType

	rfq.Sign = makeSign(param, Key)

	xmlParam := paramToXML(rfq)

	// 执行http请求
	response, e := sendXmlRequest("POST", WxHost+"/pay/closeorder", xmlParam)
	if e != nil {
		return
	}
	// 解析xml
	e = xml.Unmarshal(response, &rqr)

	if e != nil {
		return
	} else if rqr.ReturnCode != "SUCCESS" {
		e = errors.New(rqr.ReturnMsg)
		return
	}

	re = rqr
	return
}
