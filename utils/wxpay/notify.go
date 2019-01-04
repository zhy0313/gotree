package wxpay

import (
	"encoding/xml"

	"jryghq.cn/utils"
)

type Notify struct {
	XmlName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
	ResultCode string   `xml:"result_code"`  // 业务结果
	ErrCode    string   `xml:"err_code"`     // 错误码
	ErrCodeDes string   `xml:"err_code_des"` // 错误描述
	AppId      string   `xml:"appid"`        // 公众账号ID
	MchId      string   `xml:"mch_id"`       // 商户号
	DeviceInfo string   `xml:"device_info"`  // 设备号
	NonceStr   string   `xml:"nonce_str"`    // 随机字符串
	Sign       string   `xml:"sign"`         // 签名

	OpenId        string `xml:"openid"`         // 用户标识
	IsSubscribe   string `xml:"is_subscribe"`   // 是否关注公众账号
	TradeType     string `xml:"trade_type"`     // 交易类型
	BankType      string `xml:"bank_type"`      // 付款银行
	TotalFee      string `xml:"total_fee"`      // 总金额
	FeeType       string `xml:"fee_type"`       // 货币种类
	CashFee       string `xml:"cash_fee"`       // 现金支付金额
	CashFeeType   string `xml:"cash_fee_type"`  // 现金支付货币类型
	CouponFee     string `xml:"coupon_fee"`     // 代金券或立减优惠金额
	CouponCount   string `xml:"coupon_count"`   // 代金券或立减优惠使用数量
	TransactionId string `xml:"transaction_id"` // 微信支付订单号
	OutTradeNo    string `xml:"out_trade_no"`   // 商户订单号
	Attach        string `xml:"attach"`         // 商家数据包，原样返回
	TimeEnd       string `xml:"time_end"`       // 支付完成时间
}

func NewNotify(payChannel string) *Notify {
	notify := &Notify{}
	var wechatAccount string
	if payChannel == "call" {
		wechatAccount = CALL_WECHAT_ACCOUNT
	} else if payChannel == "mini" {
		wechatAccount = SP_WECHAT_ACCOUNT
	}
	// 解析微信账号
	Account, e := decodeWechatAccount(wechatAccount)
	if e != nil {
		utils.Log().WriteWarn("解析微信账号错误, payChannel:", payChannel, "wechatAccount:", wechatAccount)
	}

	Key = Account.Key
	return notify
}

func (this *Notify) VerifySignUrl(data string, sign string) (ok bool, err error) {
	ok = false
	resp := Notify{}
	err = xml.Unmarshal([]byte(data), &resp)
	if err != nil {
		return
	}
	xmlMap, err := xml2Map(resp)
	if err != nil {
		return
	}
	newSign := notifyMakeSign(xmlMap, Key)
	if newSign == sign {
		ok = true
	}
	return
}
