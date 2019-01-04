package gaode

import "math/rand"

//distance sdk
type Distance struct {
	Key string
}

//Spot 初始化
func (d *Distance) AmapKey() *Distance {
	var params []string
	params = append(params, "467d720a723dd21df7415338cbaf0c6a")
	params = append(params, "990f524650d2f8728ee6d9ab869fe0ed")
	params = append(params, "3a9ee4346ed40ecd3d70537e784311aa")
	params = append(params, "dcd4028d777ce9dc94e86ba7d40680cb") // 新 key 牛逼的很嘞

	randNum := rand.Intn(3)

	d.Key = params[randNum]
	return d
}
