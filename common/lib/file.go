package lib

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
)

var ConfEnvPath string //配置文件夹
var ConfEnv string     //配置环境名 dev/debug/test/pro

//解析配置文件目录
//配置文件需要放到一个文件夹中
//eg: config = conf/dev/base.toml

func ParseConfPath(config string) (err error) {
	path := strings.Split(config, string(os.PathSeparator))
	prefix := strings.Join(path[:len(path)-1], string(os.PathSeparator))
	ConfEnvPath = prefix
	ConfEnv = path[len(path)-2]
	return
}

func GetConfEnv() string {
	return ConfEnv
}

func GetConfEnvPath() string {
	return ConfEnvPath
}
func GetConfPath(fileName string) string {
	return ConfEnvPath + "/" + fileName + ".toml"
}

func GetConfFilePath(fileName string) string {
	return ConfEnvPath + "/" + fileName
}
func ParseConfig(path string, conf interface{}) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open config %s fail,%v", path, err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("read config fail,%v", err)
	}

	v := viper.New()
	if ConfEnv == "" {
		v.SetConfigType("toml")
	} else {
		fileType := strings.Split(ConfEnv, ".")
		v.SetConfigType(fileType[len(fileType)-1])
	}
	v.ReadConfig(bytes.NewBuffer(data))
	if err := v.Unmarshal(conf); err != nil {
		return fmt.Errorf("parse config failed,config:%s,err:%v", path, err)
	}

	return nil
}
