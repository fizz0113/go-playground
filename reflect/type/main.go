package main

import (
	"fmt"
	"reflect"
)

func inspectType(i interface{}) {
	// 獲取類型
	t := reflect.TypeOf(i)
	// 獲取值
	v := reflect.ValueOf(i)

	fmt.Printf("type: %v, value: %v\n", t, v)

	// 如果是結構體，列舉字段
	if t.Kind() == reflect.Struct {
		fmt.Println("  struct fields:")
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			fmt.Printf("    name: %v, type: %v\n", field.Name, field.Type)
		}
	}
}

func main() {
	// 測試不同類型
	inspectType(42)                    // 整數
	inspectType("hello")               // 字串
	inspectType(struct{ Name string }{ // 結構體
		Name: "Go",
	})
}
