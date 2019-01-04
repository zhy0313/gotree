package utils

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

//InArray 是否在数组内
func InArray(array interface{}, item interface{}) bool {
	values := reflect.ValueOf(array)
	if values.Kind() != reflect.Slice {
		return false
	}

	size := values.Len()
	list := make([]interface{}, size)
	slice := values.Slice(0, size)
	for index := 0; index < size; index++ {
		list[index] = slice.Index(index).Interface()
	}

	for index := 0; index < len(list); index++ {
		if list[index] == item {
			return true
		}
	}
	return false
}

//HttpBuildQuery转换get参数
func HttpBuildQuery(args map[string]interface{}) string {
	result := ""
	for k, v := range args {
		result += k + "=" + fmt.Sprint(v) + "&"
	}

	return url.PathEscape(strings.TrimSuffix(result, "&"))
}

func TakeSliceArg(arg interface{}) (out []interface{}, ok bool) {

	slice := reflect.ValueOf(arg)
	if slice.Kind() != reflect.Slice {
		return nil, false
	}

	c := slice.Len()
	out = make([]interface{}, c)
	for i := 0; i < c; i++ {
		out[i] = slice.Index(i).Interface()
	}
	return out, true
}

// 分页
func PageArray(array interface{}, index, size int) (arr []interface{}, totalPage int) {
	list, flag := TakeSliceArg(array)
	if !flag {
		return
	}
	lenth := len(list)
	pageOffSet := 0
	if lenth > 0 {
		if index > 1 {
			pageOffSet = (index - 1) * size
			if size < (lenth - pageOffSet) {
				arr = list[pageOffSet : size+pageOffSet]
			} else {
				if pageOffSet < lenth {
					arr = list[pageOffSet:]
				}
			}
		} else {
			if lenth > size {
				arr = list[0:size]
			} else {
				arr = list[0:]
			}
		}
	}
	totalPage = (lenth + size - 1) / size
	return
}

//ArrayMerge 数组合并   ArrayMerge(&dsc, src1, src2)
func ArrayMerge(dsc interface{}, src1 interface{}, src2 interface{}) error {
	return nil
}

//ArrayCopy 数组拷贝
func ArrayCopyMake(dsc interface{}, src interface{}) error {
	dscValue := reflect.ValueOf(dsc)
	if dscValue.Elem().Kind() != reflect.Slice {
		return errors.New("dsc error")
	}

	srcValue := reflect.ValueOf(src)
	if srcValue.Kind() != reflect.Slice {
		return errors.New("src error")
	}

	val := reflect.ValueOf(dsc)
	sInd := reflect.Indirect(val)
	etyp := sInd.Type().Elem()
	dsctyp := etyp
	if dsctyp.Kind() == reflect.Ptr {
		dsctyp = dsctyp.Elem()
	}

	val = reflect.ValueOf(src)
	sInd = reflect.Indirect(val)
	etyp = sInd.Type().Elem()
	srctyp := etyp
	if srctyp.Kind() == reflect.Ptr {
		srctyp = srctyp.Elem()
	}

	result := reflect.MakeSlice(reflect.TypeOf(dsc).Elem(), srcValue.Len(), srcValue.Len())
	structMode := false
	if dsctyp.Kind() != srctyp.Kind() {
		return errors.New("非法")
	}

	if dsctyp.Kind() == reflect.Struct {
		structMode = true
	}

	for index := 0; index < srcValue.Len(); index++ {
		if !structMode {
			result.Index(index).Set(srcValue.Index(index))
			continue
		}
		for j := 0; j < dsctyp.NumField(); j++ {
			v := srcValue.Index(index).FieldByName(dsctyp.Field(j).Name)
			if !v.IsValid() {
				continue
			}
			result.Index(index).FieldByName(dsctyp.Field(j).Name).Set(v)
		}
	}

	dscValue.Elem().Set(result)
	return nil
}

//ArrayNew 创建数组
func ArrayNew(dsc interface{}, len int) error {
	dscValue := reflect.ValueOf(dsc)
	if dscValue.Elem().Kind() != reflect.Slice {
		return errors.New("dsc error")
	}

	result := reflect.MakeSlice(reflect.TypeOf(dsc).Elem(), len, len)
	dscValue.Elem().Set(result)
	return nil
}

// GetSlots 生成","分隔的"?"字符串
// 例如 GetSlots(3) 转换为 "?,?,?"
func GetSlots(count int) (slots string) {
	slots = strings.Repeat("?,", count)
	slots = strings.TrimRight(slots, ",")
	return
}

//ArrayDelete 删除数组指定下标元素
func ArrayDelete(arr interface{}, indexArr ...int) {
	dscValue := reflect.ValueOf(arr)
	if dscValue.Elem().Kind() != reflect.Slice {
		Log().WriteError("dsc error")
	}

	if dscValue.Elem().Len() > 512 {
		Log().WriteError("数组过长,请使用list")
	}

	result := reflect.MakeSlice(reflect.TypeOf(arr).Elem(), 0, dscValue.Elem().Len()-len(indexArr))
	for index := 0; index < dscValue.Elem().Len(); index++ {
		if InArray(indexArr, index) {
			continue
		}
		result = reflect.Append(result, dscValue.Elem().Index(index))
	}

	dscValue.Elem().Set(result)
	return
}

//数组去重来源网络 string
func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i >= 1 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

//数组去重来源网络 int
func RemoveDuplicatesAndEmptyInt(a []int) (ret []int) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i >= 1 && a[i-1] == a[i]) || a[i] == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

type sortSlice []struct {
	data reflect.Value
	x    int
}

func (self sortSlice) Len() int {
	return len(self)
}
func (self sortSlice) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}
func (self sortSlice) Less(i, j int) bool {
	return self[j].x < self[i].x
}

//ArraySort 降序
func ArraySort(array interface{}, field string, reverse ...bool) {
	srcValue := reflect.ValueOf(array)
	if srcValue.Elem().Kind() != reflect.Slice {
		panic("ArraySort 不是切片数据或为空")
	}

	if srcValue.Elem().Len() == 0 {
		return
	}

	sortArray := make(sortSlice, srcValue.Elem().Len())
	for index := 0; index < srcValue.Elem().Len(); index++ {
		sortFiled := srcValue.Elem().Index(index).FieldByName(field)
		if !sortFiled.IsValid() {
			panic("ArraySort field不存在成员:" + field)
		}
		numFiled, err := strconv.Atoi(fmt.Sprint(sortFiled.Interface()))
		if err != nil {
			panic("ArraySort field获取int失败" + field)
		}
		sortArray[index].x = numFiled
		sortArray[index].data = srcValue.Elem().Index(index)
	}
	if len(reverse) > 0 {
		sort.Sort(sort.Reverse(sortArray))
	} else {
		sort.Sort(sortArray)
	}

	result := reflect.MakeSlice(reflect.TypeOf(array).Elem(), 0, srcValue.Elem().Len())
	for index := 0; index < len(sortArray); index++ {
		result = reflect.Append(result, sortArray[index].data)
	}
	srcValue.Elem().Set(result)
	return
}

//ArraySortReverse 升序
func ArraySortReverse(array interface{}, field string) {
	ArraySort(array, field, true)
}
