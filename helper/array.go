// Copyright gotree Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helper

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
)

//Slice 是否在数组内
func InSlice(array interface{}, item interface{}) bool {
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

//NewSlice 创建数组
func NewSlice(dsc interface{}, len int) error {
	dscValue := reflect.ValueOf(dsc)
	if dscValue.Elem().Kind() != reflect.Slice {
		return errors.New("dsc error")
	}

	result := reflect.MakeSlice(reflect.TypeOf(dsc).Elem(), len, len)
	dscValue.Elem().Set(result)
	return nil
}

//SliceDelete 删除数组指定下标元素
func SliceDelete(arr interface{}, indexArr ...int) {
	dscValue := reflect.ValueOf(arr)
	if dscValue.Elem().Kind() != reflect.Slice {
		Log().WriteError("dsc error")
	}
	result := reflect.MakeSlice(reflect.TypeOf(arr).Elem(), 0, dscValue.Elem().Len()-len(indexArr))
	for index := 0; index < dscValue.Elem().Len(); index++ {
		if InSlice(indexArr, index) {
			continue
		}
		result = reflect.Append(result, dscValue.Elem().Index(index))
	}

	dscValue.Elem().Set(result)
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

//SliceSort 降序
func SliceSort(array interface{}, field string, reverse ...bool) {
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

//SliceSortReverse 升序
func SliceSortReverse(array interface{}, field string) {
	SliceSort(array, field, true)
}
