package nixie_cloud_storage

import (
	"fmt"
	"net/http"
)

func main() {
	// 建立路由规则
	http.HandleFunc("/file/upload", UploadHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Fail to start server, err: %s", err.Error())
	}
}
