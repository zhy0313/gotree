package wxpay

import (
	"encoding/xml"
	"errors"
	"fmt"

	"jryghq.cn/utils"
)

type UnifiedOrder struct {
	WeChatClient
	TotalFee       int    `xml:"total_fee"`
	Body           string `xml:"body"`
	SpbillCreateIp string `xml:"spbill_create_ip"`
	ProductId      string `xml:"product_id"`
	OpenId         string `xml:"openid"`
}

// 初始化微信统一下单固定参数
func NewWeChatUnifiedOrder(host, payChannel, notifyUrl, tradeType, orderId, outTradeNo, nonceStr, subject string, totalFee float64, openId string) (*UnifiedOrder, error) {

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
	uo := &UnifiedOrder{}
	// 绑定配置参数
	WxHost = host
	Key = Account.Key
	// 绑定请求参数
	uo.AppId = Account.AppId
	uo.MchId = Account.MchId
	uo.SpbillCreateIp = utils.LocalIp()
	uo.SignType = "MD5" //签名类型，默认为MD5
	uo.Body = subject   //商品描述 todo根据不同支付场景定义不同
	uo.TotalFee = utils.Float64ToInt(totalFee * 100)
	uo.NotifyUrl = notifyUrl
	uo.TradeType = tradeType // APP:APP支付 NATIVE:二维码支付 JSAPI:公众号支付
	uo.OutTradeNo = outTradeNo
	uo.NonceStr = nonceStr
	uo.ProductId = orderId
	uo.OpenId = openId
	return uo, nil
}

// 微信统一下单接口
func (uo *UnifiedOrder) WeChatUnifiedOrder() (re UnifieldOrderResponse, e error) {

	uor := UnifieldOrderResponse{}

	// 组合签名参数
	param := make(map[string]string, 0)
	param["appid"] = uo.AppId
	param["mch_id"] = uo.MchId
	param["body"] = uo.Body
	param["total_fee"] = fmt.Sprintf("%d", uo.TotalFee)
	param["spbill_create_ip"] = uo.SpbillCreateIp
	param["notify_url"] = uo.NotifyUrl
	param["trade_type"] = uo.TradeType
	param["out_trade_no"] = uo.OutTradeNo
	param["nonce_str"] = uo.NonceStr
	param["sign_type"] = uo.SignType
	param["product_id"] = uo.ProductId
	param["openid"] = uo.OpenId

	uo.Sign = makeSign(param, Key)

	xmlParam := paramToXML(uo)

	// 执行http请求
	response, e := sendXmlRequest("POST", WxHost+"/pay/unifiedorder", xmlParam)
	if e != nil {
		return
	}
	// 解析xml
	e = xml.Unmarshal(response, &uor)

	if e != nil {
		return
	} else if uor.ReturnCode != "SUCCESS" {
		e = errors.New(uor.ReturnMsg)
		return
	}
	uor.OutTradeNo = uo.OutTradeNo
	re = uor
	return
}
