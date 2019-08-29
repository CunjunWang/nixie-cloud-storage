/*
 * create by CunjunWang
 * on 2019-08-28
 */

package meta

// FileMeta: 文件元信息结构体
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta: 更新全局存储map, 新增或者更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// GetFileMeta: 通过文件Sha1获取文件元信息
func GetFileMeta (fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}