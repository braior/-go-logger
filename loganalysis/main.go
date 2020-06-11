package main

import (
	"fmt"
	"sort"
)

// IntMapOrderedItem 定义结构体
type IntMapOrderedItem struct {
	Key   string
	Value int
}

// IntMapOrdered 类型为IntMapOrderedItem指针切片
type IntMapOrdered []*IntMapOrderedItem

// NewIntMapOrdered 初始化IntMapOrdered
func NewIntMapOrdered(m map[string]int) IntMapOrdered {

	// 初始化一个IntMapOrdered切片，并指定容量
	im := make(IntMapOrdered, len(m))
	index := 0
	for key, value := range m {
		//将map中的key,value打包成IntMapOrderedItem结构体赋值给im切片
		im[index] = &IntMapOrderedItem{key, value}
		index++
	}
	return im
}

// Get 获取IntMapOrdered切片里结构体的key,value
func (m IntMapOrdered) Get(key string) (int, error) {
	for _, item := range m {
		if item.Key == key {
			return item.Value, nil
		}
	}
	return 0, fmt.Errorf("key not found")
}

// Set 修改IntMapOrdered切片里结构体的key,value
func (m IntMapOrdered) Set(key string, value int) {
	for _, item := range m {
		if item.Key == key {
			item.Value = value
			return
		}
	}
	m = append(m, &IntMapOrderedItem{key, value})
}

// Delete 删除IntMapOrdered切片结构体中的key
func (m IntMapOrdered) Delete(key string) {
	index := -1
	for idx, item := range m {
		if item.Key == key {
			index = idx
			break
		}
	}

	if index != -1 {
		m = append(m[:index], m[index+1:]...)
	}
}

// Len 返回IntMapOrdered切片的长度
func (m IntMapOrdered) Len() int {
	return len(m)
}

// Less 比较IntMapOrdered切片里的结构体中指定key对应value大小
func (m IntMapOrdered) Less(i, j int) bool {
	return m[i].Value < m[j].Value
}

// Swap 交换IntMapOrdered切片里的结构体中指定key,value
func (m IntMapOrdered) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// SortMap 对IntMapOrdered切片进行排序
func SortMap(m IntMapOrdered, order string) IntMapOrdered {
	if order == "desc" {
		sort.Sort(sort.Reverse(m))
	} else {
		sort.Sort(m)
	}
	return m
}
