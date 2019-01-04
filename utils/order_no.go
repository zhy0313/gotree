package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//OrderNoYear 通过单号获取年
func OrderNoExtractYear(OrderNo string) (int, error) {
	buf := []rune(OrderNo)
	if len(buf) < 3 {
		return 0, errors.New("error OrderNoExtractOrderId OrderNo :" + fmt.Sprint(OrderNo))
	}
	i, err := strconv.ParseInt(string(buf[0:2]), 10, 32)
	if err != nil {
		return 0, errors.New("error OrderNoExtractOrderId OrderNo :" + fmt.Sprint(OrderNo))
	}
	return int(i + 2000), nil
}

//OrderNoYear 通过单号获取月
func OrderNoExtractMonth(OrderNo string) (int, error) {
	buf := []rune(OrderNo)
	if len(buf) < 3 {
		return 0, errors.New("error OrderNoExtractOrderId OrderNo :" + fmt.Sprint(OrderNo))
	}
	i, err := strconv.ParseInt(string(buf[2]), 16, 32)
	if err != nil {
		return 0, errors.New("error OrderNoExtractOrderId OrderNo :" + fmt.Sprint(OrderNo))
	}
	return int(i), nil
}

//OrderNoGeo 通过单号获经纬度
func OrderNoExtractGeo(OrderNo string) string {
	buf := []rune(OrderNo)
	if len(buf) < 5 {
		return ""
	}
	orderGeo := string(buf[len(buf)-5:])
	return strings.ToLower(orderGeo)
}

//OrderNoExtractOrderId 通过单号获id
func OrderNoExtractOrderId(OrderNo string) (int64, error) {
	buf := []rune(OrderNo)
	if len(buf) < 9 {
		return 0, errors.New("error OrderNoExtractOrderId OrderNo :" + fmt.Sprint(OrderNo))
	}
	orderId := string(buf[3 : len(buf)-5])
	i, err := strconv.ParseInt(orderId, 36, 64)
	if err != nil {
		return 0, errors.New("error OrderNoExtractOrderId OrderNo :" + fmt.Sprint(OrderNo))
	}
	return i, nil
}

//OrderNoCreate 创建订单号
func OrderNoCreate(latitude, longitude float64, orderId int64) string {
	geo := GeoHashEncode(latitude, longitude, 5)
	year := time.Now().Year() - 2000
	month := fmt.Sprintf("%x", int(time.Now().Month()))
	orderIDx := strconv.FormatInt(orderId, 36)
	orderNO := strings.ToUpper(fmt.Sprint(year, month, orderIDx, geo))
	return orderNO
}

//InvoiceCode 创建发票号
func InvoiceCode(id int64) string {
	year := time.Now().Year() - 2000
	month := strconv.FormatInt(int64(time.Now().Month()), 16)
	orderIDx := strconv.FormatInt(id, 32)
	return strings.ToUpper(fmt.Sprint(year, month, orderIDx))
}

//InvoiceCodeExtractId 提取发票 id
func InvoiceCodeExtractId(code string) (id int64, e error) {
	buf := []rune(code)
	if len(buf) < 3 {
		return 0, errors.New("error InvoiceCodeExtractId code :" + fmt.Sprint(code))
	}
	idx := string(buf[3:])
	i, err := strconv.ParseInt(idx, 32, 64)
	if err != nil {
		return 0, errors.New("error InvoiceCodeExtractId code :" + fmt.Sprint(code))
	}
	id = i
	return
}
