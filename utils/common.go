package utils

// // 生成流水单号
// func CreateTradeNo(prefix string) string {
// 	return fmt.Sprintf("%s%s%d", prefix, time.Now().Format("20060102150405"), randInt(10000, 99999))
// }

// //RandInt 生成随机数
// func randInt(min int, max int) int {
// 	if max-min <= 0 {
// 		return min
// 	}
// 	rand.Seed(time.Now().UTC().UnixNano())
// 	return min + rand.Intn(max-min)
// }
