package config

import (
	"github.com/alyu/configparser"
	"github.com/shellus/pkg/logs"
	"github.com/shellus/pkg/sutil"
	"path/filepath"
)

type ossConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	Bucket          string
}

func GetConfig() *ossConfig {
	configPath := filepath.Join(sutil.HomeDir(), ".ossutilconfig")
	if !sutil.FileExists(configPath) {
		logs.Fatal("configfile: %s not found", configPath)
	}
	config, err := configparser.Read(configPath)
	if err != nil {
		logs.Fatal("read configfile err: ",err)
	}

	section, err := config.Section("Credentials")
	if err != nil {
		logs.Fatal("config Section Credentials not found : %s", err)
	}

	if !section.Exists("endpoint") {
		logs.Fatal("config key endpoint not found")
	}
	if !section.Exists("accessKeyID") {
		logs.Fatal("config key accessKeyID not found")
	}
	if !section.Exists("accessKeySecret") {
		logs.Fatal("config key accessKeySecret not found")
	}
	if !section.Exists("bucket") {
		logs.Fatal("config key bucket not found")
	}
	c := new(ossConfig)
	c.Endpoint = section.ValueOf("endpoint")
	c.AccessKeyID = section.ValueOf("accessKeyID")
	c.AccessKeySecret = section.ValueOf("accessKeySecret")
	c.Bucket = section.ValueOf("bucket")
	return c
}
