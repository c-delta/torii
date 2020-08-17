package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

// Settings 設定檔結構
type Settings struct {
	// TCP Server
	Host string           `yaml:"host"`
	SSL  serverSSLOptions `yaml:"ssl"`
	Port int32            `yaml:"port"`

	//
	Name          string `yaml:"name"`
	Version       string `yaml:"version"`
	MaxPacketSize int32  `yaml:"max_packet_size"`

	// // Config
	// ConfFilePath string

	// Log
	LogFile    string `yaml:"log_file"`
	DebugLevel string `yaml:"debug_level"`
}

type serverSSLOptions struct {
	Enable bool   `yaml:"enable"`
	Cert   string `yaml:"cert"`
	Key    string `yaml:"key"`
	Pem    string `yaml:"pem"`
}

// Config 匯出用
func Config() *Settings {
	conf := &Settings{}
	conf.Reload()
	return conf
}

func (conf *Settings) translateConfig(ConfFilePath string) error {
	data, err := ioutil.ReadFile(ConfFilePath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		panic(err)
	}
	return err
}

// IsPathExists 判斷路徑是否存在
func IsPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CurrentPath 目前路徑
func CurrentPath() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return path
}

// Reload 重新載入設定檔
func (conf *Settings) Reload() {
	ConfFilePath := path.Join(CurrentPath(), "torii.yaml")
	if isExists, _ := IsPathExists(ConfFilePath); isExists == false {
		panic(errors.New("No Config YAML"))
	}

	err := conf.translateConfig(ConfFilePath)
	if err != nil {
		panic(err)
	}
}
