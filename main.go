package main

import (
	"fmt"
	"net/http"
	"nixie-cloud-storage/handler"
)

func main() {
	// 建立路由规则
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/success", handler.UploadSuccessHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Fail to start server, err: %s", err.Error())
	}
}
