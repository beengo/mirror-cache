package main

import (
	"github.com/beengo/mirror-cache/logs"
	"github.com/jinzhu/configor"
)

// Mirror mirror config
type Mirror struct {
	Name        string
	Prefix      string `default:"/"`
	ProxyTarget string
	LocalDir    string
	Expire      int64 `default:"3600*24"`
}

// ConfigStruct config struct
type ConfigStruct struct {
	Listen    string `default:":8080"`
	BlockSize int64  `default:"131072"`
	Mirrors   []Mirror
	Logger    logs.LogHandler
}

// Config config
var Config ConfigStruct

// LoadConfig load config
func LoadConfig(filePath string) {
	configor.Load(&Config, filePath)
}
