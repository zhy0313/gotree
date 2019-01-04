package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"reflect"
	"time"

	"net/mail"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"

	mrand "math/rand"

	"github.com/astaxie/beego"
	"github.com/zheng-ji/goSnowFlake"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

//Left 字串截取
func Left(s string, length int) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	if len(runes) > length {
		return string(runes[0:length]) + ".."
	}
	return s
}
func GetFileSuffix(s string) string {
	re, _ := regexp.Compile(".(jpg|jpeg|png|gif|exe|doc|docx|ppt|pptx|xls|xlsx)")
	suffix := re.ReplaceAllString(s, "")
	return suffix
}

func RandInt64(min, max int64) int64 {
	maxBigInt := big.NewInt(max)
	i, _ := rand.Int(rand.Reader, maxBigInt)
	if i.Int64() < min {
		return RandInt64(min, max)
	}
	return i.Int64()
}

func Strim(str string) string {
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	return str
}

func Unicode(rs string) string {
	json := ""
	for _, r := range rs {
		rint := int(r)
		if rint < 128 {
			json += string(r)
		} else {
			json += "\\u" + strconv.FormatInt(int64(rint), 16)
		}
	}
	return json
}

func HTMLEncode(rs string) string {
	html := ""
	for _, r := range rs {
		html += "&#" + strconv.Itoa(int(r)) + ";"
	}
	return html
}

/**
 *  to: example@example.com;example1@163.com;example2@sina.com.cn;...
 *  subject:The subject of mail
 *  body: The content of mail
 */
func SendMail(to string, subject string, body string) error {
	user := beego.AppConfig.String("mailfrom")
	password := beego.AppConfig.String("mailpassword")
	host := beego.AppConfig.String("mailhost")

	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	content_type = "Content-type:text/html;charset=utf-8"

	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

var idWorker *goSnowFlake.IdWorker

func SnowFlakeID() int64 {
	if idWorker == nil {
		mrand.Seed(time.Now().UnixNano())
		machineID := int64(1 + mrand.Intn(980))
		idWorker, _ = goSnowFlake.NewIdWorker(machineID)
	}
	if id, err := idWorker.NextId(); err != nil {
		return 0
	} else {
		return id
	}
}

// func Iconv(str string, srcCode string, tagCode string) (result string, err error) {
// 	cd, err := iconv.Open(tagCode, srcCode)
// 	if err != nil {
// 		return
// 	}
// 	defer cd.Close()
// 	result = cd.ConvString(str)
// 	return
// }
//Utf8ToGbk utf8编码转gbk
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

//GbkToUtf8 gbk编码转utf8
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// func BytesConvert(bytes []byte, reply interface{}) error {
// 	// logs.Info("获取类型获取类型：", reply)
// 	switch reply.(type) {
// 	case *string:
// 		b, err := Encode(string(bytes))
// 		err = Decode(b, reply)
// 		return err
// 	case *float32:
// 		var s string
// 		s = string(bytes)
// 		result, err := strconv.ParseFloat(s, 32) //32/64
// 		b, err := Encode(result)
// 		err = Decode(b, reply)
// 		return err
// 	case *float64:
// 		var s string
// 		s = string(bytes)
// 		result, err := strconv.ParseFloat(s, 64) //32/64
// 		b, err := Encode(result)
// 		err = Decode(b, reply)
// 		return err
// 	case *int:
// 		var s string
// 		s = string(bytes)
// 		result, err := strconv.Atoi(s)
// 		b, err := Encode(result)
// 		err = Decode(b, reply)
// 		return err
// 	case *int8:
// 		var s string
// 		s = string(bytes)
// 		result, err := strconv.ParseInt(s, 10, 8) //值，基数，类型
// 		b, err := Encode(result)
// 		err = Decode(b, reply)
// 		return err
// 	case *int32:
// 		var s string
// 		s = string(bytes)
// 		result, err := strconv.ParseInt(s, 10, 32)
// 		b, err := Encode(result)
// 		err = Decode(b, reply)
// 		return err
// 	case *int64:
// 		var s string
// 		s = string(bytes)
// 		result, err := strconv.ParseInt(s, 10, 64)
// 		b, err := Encode(result)
// 		err = Decode(b, reply)
// 		return err
// 	default: //struct
// 		var s string
// 		var n int64
// 		s = string(bytes)
// 		n = int64(len(s)) + 10 //获取字节长度
// 		headerReader := io.LimitReader(strings.NewReader(s), n)
// 		err := json.NewDecoder(headerReader).Decode(reply)
// 		return err

// 		// return []byte(reply), nil
// 		// case nil:
// 		// 	return nil, ErrNil
// 		// case Error:
// 		// 	return nil, reply
// 	}
// 	logs.Error("BytesConvert转换异常", reply)
// 	return errors.New("BytesConvert转换异常")
// 	//集中基本类型转换
// 	//http://www.jb51.net/article/119164.htm
// 	//http://blog.csdn.net/pangudashu/article/details/50695924
// 	//json解析
// 	//http://blog.csdn.net/shachao888/article/details/53840577

// }

const (
	regular = "^1\\d{10}$"
)

func Validate(mobileNum string) bool {
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func HideMobile(mobileNum string) (HideMobileNum string) {
	if len(mobileNum) < 11 {
		return mobileNum
	}
	HideMobileNum = Substr(mobileNum, 0, 3) + fmt.Sprint("*****", Substr(mobileNum, len(mobileNum)-3, 3))
	return
}

func RegMatch(str string, reg string) bool {
	r := regexp.MustCompile(reg)
	return r.MatchString(str)
}
func RegFindString(str string, key string, reg string) (result string) {
	r := regexp.MustCompile(reg)
	m := r.FindStringSubmatch(str)
	if m != nil {
		for i, name := range r.SubexpNames() {
			if name == key {
				return m[i]
			}
		}
	}
	return
}

// ValidateBankCardID 银行卡号验证
func ValidateBankCardID(cardId string) bool {
	var oddSum int
	var evenSum int

	oddSum = 0
	evenSum = 0

	cardId = reverseString(cardId)

	for i, _ := range cardId {

		item, _ := strconv.Atoi(string(cardId[i]))
		//fmt.Println("", item)

		if i%2 == 1 {
			num := item * 2
			sum := num
			if num > 9 {
				sum = 0
				for _, n := range strconv.Itoa(num) {
					inum, _ := strconv.Atoi(string(n))
					sum += inum
				}
			}
			evenSum += sum
		} else {
			oddSum += item
		}
		//fmt.Println(oddSum, "-----", evenSum)
		//fmt.Println("--------------")
	}

	//fmt.Println("", oddSum, evenSum)
	return (oddSum+evenSum)%10 == 0
}

// reverseString字符串反转
func reverseString(s string) string {
	runes := []rune(s)

	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}

	return string(runes)
}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

// FromIntArray 将[]int数组转为","分隔的字符串
// 例如 [1,2,3] 转换为 "1,2,3"
func FromIntArray(arr []int) string {
	n := len(arr)
	strArr := make([]string, n)
	for i := 0; i < n; i++ {
		strArr[i] = strconv.Itoa(arr[i])
	}
	return strings.Join(strArr, ",")
}

// FromStringArray 将数组转为","分隔的字符串
// 例如 [1,2,3] 转换为 "1,2,3"
func FromStringArray(arr []string) string {
	n := len(arr)
	strArr := make([]string, n)
	for i := 0; i < n; i++ {
		strArr[i] = arr[i]
	}
	return strings.Join(strArr, ",")
}

// FromInt64Array 将[]int数组转为","分隔的字符串
// 例如 [1,2,3] 转换为 "1,2,3"
func FromInt64Array(arr []int64) string {
	n := len(arr)
	strArr := make([]string, n)
	for i := 0; i < n; i++ {
		strArr[i] = strconv.FormatInt(arr[i], 10)
	}
	return strings.Join(strArr, ",")
}

// GetIntArray 将","分隔的字符串转换为数字数组
// 例如 "1,2,3" 转换为 [1, 2, 3]
func GetIntArray(str string) (ret []int, err error) {
	for _, item := range strings.Split(str, ",") {
		if item != "" {
			var value int
			value, err = strconv.Atoi(item)
			if err != nil {
				return
			}
			ret = append(ret, value)
		}
	}
	return
}

// GetStringArray 将","分隔的字符串转换为字符串数组
// 例如 "1,2,3" 转换为 ["1", "2", "3"]
func GetStringArray(str string) (ret []string) {
	for _, item := range strings.Split(str, ",") {
		if item != "" {
			ret = append(ret, item)
		}
	}
	return
}

// GetfloatArray 将","分隔的字符串转换为数字数组
// 例如 "1.1,2,3" 转换为 [1.1, 2, 3]
func Getfloat64Array(str string) (ret []float64, err error) {
	for _, item := range strings.Split(str, ",") {
		if item != "" {
			var value float64
			value, err = strconv.ParseFloat(item, 64)
			if err != nil {
				return
			}
			ret = append(ret, value)
		}
	}
	return
}

// GetInt64Array 将","分隔的字符串转换为数字数组
// 例如 "1,2,3" 转换为 [1, 2, 3]
func GetInt64Array(str string) (ret []int64, err error) {
	for _, item := range strings.Split(str, ",") {
		if item != "" {
			var value int64
			value, err = strconv.ParseInt(item, 10, 64)
			if err != nil {
				return
			}
			ret = append(ret, value)
		}
	}
	return
}

//ReplaceSpecialCharacters 过滤特殊字符
func ReplaceSpecialCharacters(str string) string {
	if len(str) > 0 {
		reg := regexp.MustCompile("[';%]")
		str = reg.ReplaceAllString(str, "")
	}
	return str
}

//RegReplace 正则替换
func RegReplace(str string, reg string) string {
	if len(str) > 0 {
		reg := regexp.MustCompile(reg)
		str = reg.ReplaceAllString(str, "")
	}
	return str
}

// dial using TLS/SSL
func dial(addr string) (*tls.Conn, error) {
	return tls.Dial("tcp", addr, nil)
}

// compose message according to "from, to, subject, body"
func composeMsg(from string, to string, subject string, body string) (message string) {
	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["Content-type"] = "text/html;charset=utf-8"
	// Setup message
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body
	return
}

// send email over SSL
func SendEmailSSL(toAddr string, subject string, body string) (err error) {
	username := beego.AppConfig.String("mailfrom")
	password := beego.AppConfig.String("mailpassword")
	servername := beego.AppConfig.String("mailhost")
	host, _, _ := net.SplitHostPort(servername)
	// get SSL connection
	conn, err := dial(servername)
	if err != nil {
		return
	}
	// create new SMTP client
	smtpClient, err := smtp.NewClient(conn, host)
	if err != nil {
		return
	}
	// Set up authentication information.
	auth := smtp.PlainAuth("", username, password, host)
	// auth the smtp client
	err = smtpClient.Auth(auth)
	if err != nil {
		return
	}
	// set To && From address, note that from address must be same as authorization user.
	from := mail.Address{"", username}
	to := mail.Address{"", toAddr}
	err = smtpClient.Mail(from.Address)
	if err != nil {
		return
	}
	err = smtpClient.Rcpt(to.Address)
	if err != nil {
		return
	}
	// Get the writer from SMTP client
	writer, err := smtpClient.Data()
	if err != nil {
		return
	}
	// compose message body
	message := composeMsg(from.String(), to.String(), subject, body)
	// write message to recp
	_, err = writer.Write([]byte(message))
	if err != nil {
		return
	}
	// close the writer
	err = writer.Close()
	if err != nil {
		return
	}
	// Quit sends the QUIT command and closes the connection to the server.
	smtpClient.Quit()
	return nil
}

//对比相同struct下不同的值，修改日志专用，返回map类型
func Contrast(DataOld, DataNew interface{}) map[string]interface{} {

	m := make(map[string]interface{})
	old := reflect.TypeOf(DataOld)
	new := reflect.TypeOf(DataNew)
	if old != new {
		return m
	}
	oldVal := reflect.ValueOf(DataOld)
	newVal := reflect.ValueOf(DataNew)
	nums := old.NumField()
	for i := 0; i < nums; i++ {
		if old.Field(i).Name == new.Field(i).Name && oldVal.Field(i).Interface() != newVal.Field(i).Interface() {
			//m[old.Field(i).Name] = fmt.Sprintf("%T -> %T", oldVal.Field(i).Interface(), newVal.Field(i).Interface())
			switch reflect.TypeOf(oldVal.Field(i).Interface()).String() {
			//多选语句switch
			case "string":
				m[old.Field(i).Name] = fmt.Sprintf("%s -> %s", oldVal.Field(i).Interface(), newVal.Field(i).Interface())
			case "int":
				m[old.Field(i).Name] = fmt.Sprintf("%d -> %d", oldVal.Field(i).Interface(), newVal.Field(i).Interface())
			case "float64":
				m[old.Field(i).Name] = fmt.Sprintf("%f -> %f", oldVal.Field(i).Interface(), newVal.Field(i).Interface())
			}

		}
	}
	return m

}
