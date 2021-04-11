package main

import "net/http"

func main() {
	// 靜態檔案（例如 index.html、js、css）
	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/", fs)

	http.ListenAndServe(":8888", nil)
}
