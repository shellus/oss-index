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
func getAllPath(bucket *oss.Bucket)map[string][]Object{

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

	prefixs := make(map[string][]Object)

	// 防止根目录下无文件
	prefixs["/"] = []Object{}

	for _, i := range objects {

		pathParts := strings.Split(i.Key, "/")

		// 分割成路径片段后，最后一个成员要么是空，要么是文件名
		// 无论如何，去掉最后一个，得到所在目录、然后将当前path加入到这个目录
		pathParts = pathParts[:len(pathParts) - 1]
		dir := strings.Join(pathParts, "/") + "/"

		prefixs[dir] = append(prefixs[dir], Object{Key:i.Key, IsDir: false, Size: i.Size})
	}

	// 把下级目录作为上级目录下的一个object
	for k := range prefixs {


		pathParts := strings.Split(k, "/")
		pathParts = pathParts[:len(pathParts) - 2]
		// 得到上级目录
		dir := strings.Join(pathParts, "/") + "/"

		// 如果上级目录存在，则把当前目录加入作为object
		if _, ok := prefixs[dir]; ok {
			prefixs[dir] = append(prefixs[dir], Object{Key: k, IsDir: true, Size: 0})
		}
	}
	return prefixs
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