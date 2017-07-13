package oss_index

import ()
import (
	"sync"
	"github.com/shellus/pkg/logs"
)

// 元文件名
const metaFileName string = ".oss_index_meta"
// 一个文件
type Object struct {
	Title string `json:"title"`
	Key   string `json:"key"`
	Size  int64 `json:"size"`
	IsDir bool `json:"is_dir"`
}
// 一个页面的json
type PathMeta struct {
	Prefix  string `json:"prefix"`
	Objects []Object `json:"objects"`
}
var prefixChan = make(chan *PathMeta)

var prefixWait = sync.WaitGroup{}

func Main() {
	bucket := getBucket()
	prefixs := getAllPath(bucket)

	// 启动10个线程
	for i := 0; i < 20; i++ {
		go thrUpdateMetaInfo()
	}

	logs.Info("prefixs count %d", prefixs)
	// 分配任务
	for k, v := range prefixs {
		prefixWait.Add(1)
		prefixChan <- &PathMeta{Prefix:k, Objects: v}
	}

	// 等待任务完成
	prefixWait.Done()
}

func thrUpdateMetaInfo(){
	bucket := getBucket()
	for prefix := range prefixChan{
		updateMetaInfo(bucket, prefix)
		prefixWait.Done()
	}
}


