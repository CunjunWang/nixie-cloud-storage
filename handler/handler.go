package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"nixie-cloud-storage/meta"
	"nixie-cloud-storage/util"
	"os"
	"time"
)

const uploadPath = "/tmp/"

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
	} else {
		// 接收文件流并存储到本地目录
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data from file, err: %s", err.Error())
			return
		}
		// 在函数退出之前把file流关掉
		defer file.Close()

		var fileLocation = uploadPath + header.Filename;

		// 存储元信息
		fileMeta := meta.FileMeta{
			FileName: header.Filename,
			Location: fileLocation,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		// 创建本地文件
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to create local file, err: %s", err.Error())
			return
		}
		defer newFile.Close()

		// 把上传的文件流拷贝写入创建的本地文件
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to write data to local file, err: %s", err.Error())
			return
		}

		// 计算文件SHA1值
		newFile.Seek(0, 0) // 移动文件句柄到文件头部, 来获取整个文件的SHA1值
		fileMeta.FileSha1 = util.FileSha1(newFile)

		// 把文件元信息存储下来
		meta.UpdateFileMeta(fileMeta)

		http.Redirect(w, r, "/file/upload/success", http.StatusFound)
	}
}

// UploadSuccessHandler: 上传完成
func UploadSuccessHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}
