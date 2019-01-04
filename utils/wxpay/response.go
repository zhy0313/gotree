package wxpay

type BaseResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`  // 业务结果
	ErrCode    string `xml:"err_code"`     // 错误码
	ErrCodeDes string `xml:"err_code_des"` // 错误描述
	AppId      string `xml:"appid"`        // 公众账号ID
	MchId      string `xml:"mch_id"`       // 商户号
	NonceStr   string `xml:"nonce_str"`    // 随机字符串
	Sign       string `xml:"sign"`         // 签名
}

// 统一下单接口返回
type UnifieldOrderResponse struct {
	BaseResponse
	DeviceInfo string `xml:"device_info"` // 设备号
	PrepayId   string `xml:"prepay_id"`   // 预支付交易会话标识
	CodeUrl    string `xml:"code_url"`    // 二维码链接
	TradeType  string `xml:"trade_type"`  // 交易类型 JSAPI 公众号支付 NATIVE 扫码支付 APP APP支付
	OutTradeNo string `xml:"out_trade_no"`
}

// 退款接口返回
type RefundResponse struct {
	BaseResponse
	TransactionId string `xml:"transaction_id"` // 微信订单号
	OutTradeNo    string `xml:"out_trade_no"`   // 商户订单号
	OutRefundNo   string `xml:"out_refund_no"`  // 商户退款单号
	RefundId      string `xml:"refund_id"`      // 微信退款单号
	RefundFee     int    `xml:"refund_fee"`     // 退款金额 单位：分
	TotalFee      int    `xml:"total_fee"`      // 订单总金额 单位：分
	CashFee       int    `xml:"cash_fee"`       // 现金支付金额 单位：分
}

// 查询订单接口返回
type OrderQueryResponse struct {
	BaseResponse
	DeviceInfo     string `xml:"device_info"`      // 设备号
	OpenId         string `xml:"open_id"`          // 用户标识
	TradeType      string `xml:"trade_type"`       // 交易类型 JSAPI 公众号支付 NATIVE 扫码支付 APP APP支付
	TradeState     string `xml:"trade_state"`      // 交易状态 SUCCESS 交易成功 REFUND 转入退款  NOTPAY 未支付 CLOSED 已关闭 REVOKED 已撤销（刷卡支付） USERPAYING 用户支付中 PAYERROR 支付失败
	BankType       string `xml:"bank_type"`        // 付款银行
	CashFee        int    `xml:"cash_fee"`         // 现金支付金额 单位：分
	TransactionId  string `xml:"transaction_id"`   // 微信订单号
	TimeEnd        string `xml:"time_end"`         // 支付完成时间
	TradeStateDesc string `xml:"trade_state_desc"` // 交易状态描述

}

// 关闭订单接口返回
type OrderCloseResponse struct {
	BaseResponse
}

type RefundQueryResponse struct {
	BaseResponse
	TransactionId string `xml:"transaction_id"` // 微信订单号
	OutTradeNo    string `xml:"out_trade_no"`   // 商户订单号
	TotalFee      string `xml:"total_fee"`      // 订单总金额
	FeeType       string `xml:"fee_type"`       // 订单金额货币种类
	CashFee       string `xml:"cash_fee"`       // 现金支付金额
	RefundCount   string `xml:"refund_count"`   // 退款笔数
}
