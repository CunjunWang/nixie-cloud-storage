package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// UploadHandler: 处理文件上传
// w: 用于输出, r: 输入指针
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传html页面
		data, err := ioutil.ReadFile("./static/view/upload.html")
		if err != nil {
			io.WriteString(w, "internal server error")
			return
		}
		io.WriteString(w, string(data))
	} else {
		// 接收文件流并存储到本地目录
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data from file, err: %s", err.Error())
			return
		}
		// 在函数退出之前把file流关掉
		defer file.Close()

		// 创建本地文件
		newFile, err := os.Create("./upload/" + header.Filename)
		if err != nil {
			fmt.Printf("Failed to create local file, err: %s", err.Error())
			return
		}
		defer newFile.Close()

		// 把上传的文件流拷贝写入创建的本地文件
		_, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to write data to local file, err: %s", err.Error())
			return
		}

		http.Redirect(w, r, "/file/upload/success", http.StatusFound)
	}
}

// UploadSuccessHandler: 上传完成
func UploadSuccessHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}
