package main

import (
	"fmt"
	"reflect"
)

// 定義一個結構體
type MyStruct struct {
	Name string
}

// 為結構體定義一個方法
func (m MyStruct) Greet(greeting string) string {
	return fmt.Sprintf("%s, %s!", greeting, m.Name)
}

func callStructMethod(s interface{}, methodName string, args ...interface{}) {
	// 獲取結構體的值
	v := reflect.ValueOf(s)

	// 獲取方法
	method := v.MethodByName(methodName)
	if !method.IsValid() {
		fmt.Printf("Method %s not found\n", methodName)
		return
	}

	// 構造 reflect.Value 的參數切片
	inputs := make([]reflect.Value, len(args))
	for i, arg := range args {
		inputs[i] = reflect.ValueOf(arg)
	}

	// 呼叫方法並獲取返回值
	results := method.Call(inputs)

	// 列出返回值
	fmt.Println("Method Results:")
	for i, result := range results {
		fmt.Printf("  Result %d: %v\n", i+1, result.Interface())
	}
}

func main() {
	// 初始化結構體
	myStruct := MyStruct{Name: "Ian"}

	// 使用 reflect 呼叫結構體的方法
	callStructMethod(myStruct, "Greet", "Hello")
}
