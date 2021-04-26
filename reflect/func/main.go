package main

import (
	"fmt"
	"reflect"
)

func inspectFunc(fn interface{}) {
	t := reflect.TypeOf(fn)

	// 確認是否為函數類型
	if t.Kind() != reflect.Func {
		fmt.Println("Provided value is not a function")
		return
	}

	// 列出函數的輸入參數類型
	fmt.Println("Function Input Parameters:")
	for i := 0; i < t.NumIn(); i++ {
		fmt.Printf("  Param %d: %s\n", i+1, t.In(i))
	}

	// 列出函數的輸出參數類型
	fmt.Println("Function Output Parameters:")
	for i := 0; i < t.NumOut(); i++ {
		fmt.Printf("  Return %d: %s\n", i+1, t.Out(i))
	}
}

func inspectAndCallFunc(fn interface{}, args ...interface{}) {
	t := reflect.TypeOf(fn)

	// 確認是否為函數類型
	if t.Kind() != reflect.Func {
		fmt.Println("Provided value is not a function")
		return
	}

	// 列出函數的輸入參數類型
	fmt.Println("Function Input Parameters:")
	for i := 0; i < t.NumIn(); i++ {
		fmt.Printf("  Param %d: %s\n", i+1, t.In(i))
	}

	// 列出函數的輸出參數類型
	fmt.Println("Function Output Parameters:")
	for i := 0; i < t.NumOut(); i++ {
		fmt.Printf("  Return %d: %s\n", i+1, t.Out(i))
	}

	// 構造 reflect.Value 的參數切片
	inputs := make([]reflect.Value, len(args))
	for i, arg := range args {
		inputs[i] = reflect.ValueOf(arg)
	}

	// 呼叫函數並獲取返回值
	results := reflect.ValueOf(fn).Call(inputs)

	// 列出返回值
	fmt.Println("Function Results:")
	for i, result := range results {
		fmt.Printf("  Result %d: %v\n", i+1, result.Interface())
	}
}

func exampleFunc(a int, b string) bool {
	fmt.Printf("=== a: %d, b: %s\n", a, b)
	return len(b) > a
}

func main() {
	// 使用 reflect 檢查函數資訊
	inspectFunc(exampleFunc)

	// 使用 reflect 檢查並執行函數
	inspectAndCallFunc(exampleFunc, 3, "hello")
}
