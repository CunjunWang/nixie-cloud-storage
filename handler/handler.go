package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"nixie-cloud-storage/meta"
	"nixie-cloud-storage/util"
	"os"
	"strconv"
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

// GetFileMetaHandler: 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数
	r.ParseForm()

	filehash := r.Form["filehash"][0]
	fileMeta := meta.GetFileMeta(filehash)
	data, err := json.Marshal(fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// FileQueryHandler: 批量查询文件元信息
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	fileMetas := meta.GetLastFileMetas(limitCnt)
	data, err := json.Marshal(fileMetas)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// DownloadHandler: 文件下载
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)

	file, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=\""+fm.FileName+"\"")
	w.Write(data)
}

// FileMetaUpdateHandler: 文件元信息更新(重命名)
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden) // 403
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed) // 405
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	w.WriteHeader(http.StatusOK) // 200
	w.Write(data)
}

// FileDeleteHandler: 文件删除
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// 删除索引
	fileSha1 := r.Form.Get("filehash")
	meta.RemoveFileMeta(fileSha1)

	// 删除文件
	fMeta := meta.GetFileMeta(fileSha1)
	// TODO: 此处删除文件可能失败, 后期再处理
	os.Remove(fMeta.Location)
	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK) // 200
}