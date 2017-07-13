package oss_index

import ()

// 元文件名
const metaFileName string = ".oss_index_meta"
// 一个文件
type Object struct {
	Key   string `json:"key"`
	Size  int64 `json:"size"`
	IsDir bool `json:"is_dir"`
}
// 一个页面的json
type PathMeta struct {
	Prefix  string `json:"prefix"`
	Objects []Object `json:"objects"`
}

func Main() {
	bucket := getBucket()
	prefixs := getAllPath(bucket)
	for k, v := range prefixs {
		updateMetaInfo(bucket, &PathMeta{Prefix:k, Objects: v})
	}
}


