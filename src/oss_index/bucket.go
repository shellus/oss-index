package oss_index

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/shellus/pkg/logs"
	"github.com/shellus/oss-index/src/oss_index/config"
)
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