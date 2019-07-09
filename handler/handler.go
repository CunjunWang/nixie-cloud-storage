package handler

import (
	"io"
	"io/ioutil"
	"net/http"
)

// UploadHandler: 处理文件上传
// w: 用于输出, r: 输入指针
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internal server error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		// 接收文件流并存储到本地目录
	}

}
