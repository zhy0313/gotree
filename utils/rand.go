package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

// 生成32位MD5随机字符串
func NonceStr(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// 生成流水单号
func CreateTradeNo(prefix string) string {
	return fmt.Sprintf("%s%s%d", prefix, time.Now().Format("20060102150405"), RandInt64(10000, 99999))
}

// GetRandomString 生成字母数字随机数 length随机数位数
func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		v := rand.Intn(i + 1)
		result = append(result, bytes[r.Intn(len(bytes)-(v+length))])
	}
	return string(result)
}
