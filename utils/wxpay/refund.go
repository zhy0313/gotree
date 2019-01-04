package wxpay

// 微信退款申请接口
import (
	"encoding/xml"
	"errors"
	"fmt"
)

type Refund struct {
	WeChatClient
	TotalFee    int    `xml:"total_fee"`
	OutRefundNo string `xml:"out_refund_no"`
	RefundFee   int    `xml:"refund_fee"`
	RefundDesc  string `xml:"refund_desc"`
}

// 初始化微信退款申请固定参数
func NewRefund(host, payChannel, outTradeNo, refundReason, notifyUrl, outrefundNo, nonceStr string, totalFee, refundFee float64) (*Refund, error) {

	var wechatAccount string
	if payChannel == "call" {
		wechatAccount = CALL_WECHAT_ACCOUNT
		CertPem = CallCertPem
		KeyPem = CallKeyPem
	} else if payChannel == "mini" {
		wechatAccount = SP_WECHAT_ACCOUNT
		CertPem = SpCertPem
		KeyPem = SpKeyPem
	}

	// 解析微信账号
	Account, e := decodeWechatAccount(wechatAccount)
	if e != nil {
		e = errors.New("解析微信账号错误")
		return nil, e
	}

	rd := &Refund{}
	// 绑定配置参数
	WxHost = host
	Key = Account.Key

	rd.AppId = Account.AppId
	rd.MchId = Account.MchId
	rd.SignType = "MD5" //签名类型，默认为MD5

	rd.TotalFee = int(totalFee * 100)
	rd.RefundFee = int(refundFee * 100)
	rd.NotifyUrl = notifyUrl
	rd.OutRefundNo = outrefundNo
	rd.OutTradeNo = outTradeNo
	rd.NonceStr = nonceStr
	rd.RefundDesc = refundReason
	return rd, nil
}

// 微信退款接口
func (refd *Refund) WeChatRefund() (re RefundResponse, e error) {

	rdr := RefundResponse{}

	// 组合签名参数
	param := make(map[string]string, 0)
	param["appid"] = refd.AppId
	param["mch_id"] = refd.MchId
	param["total_fee"] = fmt.Sprintf("%d", refd.TotalFee)
	param["refund_fee"] = fmt.Sprintf("%d", refd.RefundFee)
	param["notify_url"] = refd.NotifyUrl
	param["out_trade_no"] = refd.OutTradeNo
	param["nonce_str"] = refd.NonceStr
	param["sign_type"] = refd.SignType
	param["refund_desc"] = refd.RefundDesc
	param["out_refund_no"] = refd.OutRefundNo

	refd.Sign = makeSign(param, Key)

	xmlParam := paramToXML(refd)
	// 加载证书
	tlsConfig, e := weChatTlsConfig(CertPem, KeyPem, RootPem)

	if e != nil {
		e = errors.New("微信退款使用证书错误")
		return
	}
	// 执行http请求
	response, e := sendSecureXmlRequest("POST", WxHost+"/secapi/pay/refund", xmlParam, tlsConfig, 0)
	if e != nil {
		return
	}
	// 解析xml
	e = xml.Unmarshal(response, &rdr)

	if e != nil {
		return
	} else if rdr.ReturnCode != "SUCCESS" {
		e = errors.New(rdr.ReturnMsg)
		return
	}
	rdr.OutRefundNo = refd.OutRefundNo // 返回商户退款订单号，保存日志用
	re = rdr
	return
}
