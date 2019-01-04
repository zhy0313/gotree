package wxpay

import (
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"time"
)

var (
	// 微信host 签名key
	WxHost, Key string
	// 微信证书
	CertPem, KeyPem, RootPem []byte // 本地测试时不需要 rootPem
)

type WechatAccount struct {
	AppId string
	MchId string
	Key   string
}

type WeChatClient struct {
	AppId      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	SignType   string `xml:"sign_type"`
	OutTradeNo string `xml:"out_trade_no"`
	NotifyUrl  string `xml:"notify_url"`
	TradeType  string `xml:"trade_type"`
}

// 解析微信账号
func decodeWechatAccount(wechatAccount string) (account WechatAccount, e error) {
	if wechatAccount == "" {
		return
	}
	e = json.Unmarshal([]byte(wechatAccount), &account)
	return
}

// 计算签名
func makeSign(param map[string]string, key string) string {

	sortedKeys := make([]string, 0)
	for k, _ := range param {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys) //字典排序

	var paramsString []string
	for _, k := range sortedKeys {
		v := param[k]
		if v != "" {
			paramsString = append(paramsString, fmt.Sprintf("%s=%s", k, v))
		}
	}
	paramsStr := strings.Join(paramsString, "&") //对参数按照key=value的格式，并按照参数名ASCII字典序排序

	signStr := paramsStr + "&key=" + key //拼接API密钥

	sign := md5.New()
	sign.Write([]byte(signStr))
	return strings.ToUpper(hex.EncodeToString(sign.Sum(nil))) //MD5运算，再将得到的字符串所有字符转换为大写
}

// 转为xml参数
func paramToXML(param interface{}) []byte {
	xmlOutput, err := xml.MarshalIndent(param, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return xmlOutput
}

// 发送普通Xml请求
func sendXmlRequest(method, url string, xmlParam []byte) (body []byte, err error) {

	req, err := http.NewRequest(method, url, strings.NewReader(string(xmlParam)))
	if err != nil {
		return
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	return
}

// 发送带证书Xml请求
func sendSecureXmlRequest(method, url string, xmlParam []byte, tlsConfig *tls.Config, timeout time.Duration) (body []byte, err error) {
	req, err := http.NewRequest(method, url, strings.NewReader(string(xmlParam)))
	if err != nil {
		return
	}

	client := http.Client{}

	if timeout > 0 {
		client.Timeout = timeout * time.Second
	}

	if tlsConfig != nil {
		client.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}

	resp, err := client.Do(req)
	if err != nil {
		err = errors.New("request fail")
		return
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

// 安全证书 导入顺序 cert、key、rootca
func weChatTlsConfig(certPem, keyPem, rootPem []byte) (tlsConfig *tls.Config, err error) {
	tlsConfig = new(tls.Config)

	// certPemBlock, _ := pem.Decode(certPem)
	// keyPemBlock, _ := pem.Decode(keyPem)

	// fmt.Println(keyPem)
	var cert tls.Certificate
	cert, err = tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		return
	}
	tlsConfig.Certificates = append(tlsConfig.Certificates, cert)

	if len(rootPem) > 0 {

		tlsConfig.RootCAs = x509.NewCertPool()
		tlsConfig.RootCAs.AppendCertsFromPEM(rootPem)
	}
	return
}

// xml转map
func xml2Map(in interface{}) (map[string]string, error) {
	xmlMap := make(map[string]string)

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("xml2Map only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		tagv := fi.Tag.Get("xml")

		if strings.Contains(tagv, ",") {
			tagvs := strings.Split(tagv, ",")

			switch tagvs[1] {
			case "innerXml":
				innerXmlMap, err := xml2Map(v.Field(i).Interface())
				if err != nil {
					return nil, err
				}
				for k, v := range innerXmlMap {
					if _, ok := xmlMap[k]; !ok {
						xmlMap[k] = v
					}
				}
			}
		} else if tagv != "" && tagv != "xml" {
			xmlMap[tagv] = v.Field(i).String()
		}
	}
	return xmlMap, nil
}

// 支付回调生成sign
func notifyMakeSign(params map[string]string, key string) string {
	var keys []string
	var sorted []string

	for k, v := range params {
		if k != "sign" && v != "" {
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)
	for _, k := range keys {
		sorted = append(sorted, fmt.Sprintf("%s=%s", k, params[k]))
	}

	str := strings.Join(sorted, "&")
	str += "&key=" + key

	return fmt.Sprintf("%X", md5.Sum([]byte(str)))
}
