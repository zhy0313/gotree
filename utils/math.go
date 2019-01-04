package utils

import (
	"fmt"
	"math"
	"strconv"

	"math/rand"

	"github.com/astaxie/beego"
)

//Float64Round 64位浮点数四舍五入
func Float64Round(f float64, n int) (r float64) {
	b := fmt.Sprintf("%0."+strconv.Itoa(n)+"f", f)
	r, err := strconv.ParseFloat(b, 64) //还原小数点后三位
	if err != nil {
		beego.Error("Float64Round 错误", err)
	}
	return r
}

//RoundFloat64 64位浮点数四舍五入
func RoundFloat64(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}

//RandExp 随机指数
//n : 指数梯度
//maxN :最大指数梯度
//weight :随机浮动的权重,值越高浮动越大 不可以大于1 不可小于0.01
func RandExp(n, maxN int, weight float64) float64 {
	if weight < 0.01 {
		weight = 0.01
	}
	if weight > 1 {
		weight = 1
	}
	index := float64(n)
	maxIndex := float64(maxN)

	rg := float64(rand.Intn(int(weight * 100))) // 100
	avg := weight * 100 / 2
	x := float64(rg-avg) / 100 //求浮动
	if index+x < 0 {
		x = 0
	}
	value := math.Log(index + x)
	maxValue := math.Log(maxIndex)
	if value > maxValue {
		value = 1.0
	} else if value < 0 {
		value = 0
	} else {
		value = value * (1.0 / maxValue)
	}

	return value
}

// 获取int数组里的最大值
func MaxInt(arr []int) (maxVal int) {
	if len(arr) == 0 {
		return 0
	}
	maxVal = arr[0]
	for index := 1; index < len(arr); index++ {
		if maxVal < arr[index] {
			maxVal = arr[index]
		}
	}
	return
}

// float64=>int
func Float64ToInt(f float64) (dst int) {
	b := fmt.Sprintf("%0.0f", f)
	dst, err := strconv.Atoi(b)
	if err != nil {
		beego.Error("Float64ToInt 错误", err)
	}
	return
}
