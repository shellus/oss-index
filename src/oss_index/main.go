package oss_index

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/shellus/oss-index/src/oss_index/config"
	"github.com/shellus/pkg/logs"
	"encoding/json"
	"bytes"
	"path"
	"strings"
	"context"
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

var handlePathChan = make(chan string, 100000000);

func Main() {
	handlePathChan <- ""
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动5个工作线程
	for i := 0; i < 10; i++ {
		go funHandlePathChan(ctx)
	}

	<-ctx.Done()
	if err := ctx.Err(); err != nil {
		logs.Fatal(err)
	}

}

// 生产储存桶
func getBucket() *oss.Bucket {
	c := config.GetConfig()
	ossClient, err := oss.New(c.Endpoint, c.AccessKeyID, c.AccessKeySecret)
	if err != nil {
		logs.Fatal(err)
	}
	bucket, err := ossClient.Bucket(c.Bucket)
	if err != nil {
		logs.Fatal(err)
	}
	return bucket
}

// 并行工作线程
func funHandlePathChan(parentCtx context.Context) {

	bucket := getBucket()

	for {
		select {

		case <-parentCtx.Done():
			logs.Debug("quit")

		case dir := <-handlePathChan:
			pathMeta := getPathMeta(bucket, dir)

			updateMetaInfo(bucket, pathMeta)

			// 递龟
			for _, p := range pathMeta.CommonPrefixes {
				handlePathChan <- p
			}
		}
	}
}

// 获取目录的元信息
func getPathMeta(bucket *oss.Bucket, lsPath string) *PathMeta {
	logs.Debug("List path '%s'", lsPath)
	//if exist, err := bucket.IsObjectExist(lsPath); err != nil || exist == false {
	//	if err != nil {
	//		logs.Fatal("bucket IsObjectExist %s error: %s", lsPath, err)
	//	}else {
	//		logs.Fatal("bucket object path %s not exist", lsPath)
	//	}
	//}

	result, err := bucket.ListObjects(oss.Prefix(lsPath), oss.Delimiter("/"))
	if err != nil {
		logs.Fatal(err)
	}

	index := new(PathMeta)

	index.Prefix = lsPath

	for _, i := range result.CommonPrefixes {
		index.CommonPrefixes = append(index.CommonPrefixes, i)
	}
	for _, i := range result.Objects {

		// 当前路径不注册
		if i.Key == lsPath {
			continue
		}

		// 元文件名不注册
		if len(i.Key) >= len(metaFileName) && strings.Contains(i.Key[len(i.Key) - len(metaFileName):], metaFileName) {
			continue
		}

		index.Objects = append(index.Objects, Object{Key:i.Key, Size:i.Size})
	}
	logs.Debug("List path '%s'd", lsPath)
	return index
}

// 上传元信息到路径
func updateMetaInfo(bucket *oss.Bucket, pathMeta *PathMeta) {

	jsonBuf, err := json.Marshal(pathMeta)
	if err != nil {
		logs.Fatal(err)
	}

	key := path.Join(pathMeta.Prefix, metaFileName)
	logs.Info("Wriet file '%s'", key)
	err = bucket.PutObject(key, bytes.NewReader(jsonBuf), oss.ContentType("application/json; charset=UTF-8"))
	if err != nil {
		logs.Fatal(err)
	}
	logs.Info("Wriet file '%s'd", key)
}