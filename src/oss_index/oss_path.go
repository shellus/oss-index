package oss_index

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/shellus/pkg/logs"
	"strings"
	"encoding/json"
	"path"
	"bytes"
	"io/ioutil"
	"github.com/shellus/pkg/sutil"
	"fmt"
)

func fetchAllObject(bucket *oss.Bucket) ([]oss.ObjectProperties) {
	var objects []oss.ObjectProperties

	keyMarker := ""
	for {
		options := []oss.Option{}
		if keyMarker != "" {
			options = append(options, oss.Marker(keyMarker))
		}
		options = append(options, oss.MaxKeys(1000))

		r, err := bucket.ListObjects(options...)
		if err != nil {
			logs.Fatal(err)
		}
		objects = append(objects, r.Objects...)
		logs.Debug("fetch object %d", len(r.Objects))

		if !r.IsTruncated {
			break
		}

		logs.Debug("fetch object continue...")

		keyMarker = r.NextMarker
	}
	logs.Debug("fetch object done , count %d", len(objects))
	return objects
}
func getAllPath(bucket *oss.Bucket) {

	var objects []oss.ObjectProperties

	if sutil.FileExists("oss.json") {
		buf, err := ioutil.ReadFile("oss.json")
		if err != nil {
			logs.Fatal(err)
		}

		err = json.Unmarshal(buf, &objects)
		if err != nil {
			logs.Fatal(err)
		}
	} else {
		objects = fetchAllObject(bucket)
		buf, err := json.Marshal(objects)
		if err != nil {
			logs.Fatal(err)
		}
		err = ioutil.WriteFile("oss.json", buf, 0777)
		if err != nil {
			logs.Fatal(err)
		}
	}

	//prefixs := make(map[string][]oss.ObjectProperties)

	for _, i := range objects {

		pathParts := strings.Split(i.Key, "/")

		// 分割成路径片段后，最后一个成员要么是空，要么是文件名
		// 无论如何，去掉最后一个，得到所在目录、然后将当前path加入到这个目录
		//
		fmt.Println(len(pathParts),i.Key, pathParts)

		//prefix :=
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