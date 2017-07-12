package oss_index

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/shellus/pkg/logs"
	"context"
	"time"
)

// 元文件名
const metaFileName string = ".oss_index_meta"
// 一个文件
type Object struct {
	Key  string `json:"key"`
	Size int64 `json:"size"`
}
// 一个页面的json
type PathMeta struct {
	Prefix         string `json:"prefix"`
	CommonPrefixes []string `json:"common_prefixes"`
	Objects        []Object `json:"objects"`
}

var handlePathChan = make(chan string, 10000);
var bucketChan = make(chan *oss.Bucket)
func Main() {
	handlePathChan <- ""
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for{
			go funHandlePathChan(cancel, <-bucketChan)
		}
	}()
	// 启动5个工作线程
	for i := 0; i < 20; i++ {
		bucketChan<-getBucket()
	}


	<-ctx.Done()
	if err := ctx.Err(); err != nil {
		logs.Fatal(err)
	}

}



// 并行工作线程
func funHandlePathChan(cancel context.CancelFunc, bucket *oss.Bucket) {
	TOP:
	for {
		select {

		case dir := <-handlePathChan:

			pathMeta := getPathMeta(bucket, dir)

			updateMetaInfo(bucket, pathMeta)

			// 递龟
			for _, p := range pathMeta.CommonPrefixes {
				handlePathChan <- p
			}
		case <-time.NewTimer(10*time.Second).C:
			logs.Debug("goroutine funHandlePathChan timeout")
			cancel()
			break TOP
		}
	}
}
