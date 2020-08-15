package lib

import (
	"bytes"
	"github.com/julianlee107/go-common/log"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type BaseConf struct {
	DebugMode    string    `mapstructure:"debug_mode"`
	TimeLocation string    `mapstructure:"time_location"`
	Log          LogConfig `mapstruacture:"log"`
	Base         struct {
		DebugMode    string `mapstructure:"debug_mode"`
		TimeLocation string `mapstructure:"time_location"`
	} `mapstructure:"base"`
}
type LogConfFileWriter struct {
	On              bool   `mapstructure:"on"`
	LogPath         string `mapstructure:"log_path"`
	RotateLogPath   string `mapstructure:"rotate_log_path"`
	WfLogPath       string `mapstructure:"wf_log_path"`
	RotateWfLogPath string `mapstructure:"rotate_wf_log_path"`
}

type LogConfConsoleWriter struct {
	On    bool `mapstructure:"on"`
	Color bool `mapstructure:"color"`
}

type LogConfig struct {
	Level string               `mapstructure:"log_level"`
	FW    LogConfFileWriter    `mapstructure:"file_writer"`
	CW    LogConfConsoleWriter `mapstructure:"console_writer"`
}

var (
	ConfBase     *BaseConf
	ViperConfMap map[string]*viper.Viper
)

// get base config
func GetBaseConf() *BaseConf {
	return ConfBase
}

func InitBaseConf(path string) error {
	ConfBase = &BaseConf{}
	err := ParseConfig(path, ConfBase)
	if err != nil {
		return err
	}
	if ConfBase.DebugMode == "" {
		if ConfBase.Base.DebugMode != "" {
			ConfBase.DebugMode = ConfBase.Base.DebugMode
		} else {
			ConfBase.DebugMode = "debug"
		}
	}

	if ConfBase.TimeLocation == "" {
		if ConfBase.Base.TimeLocation != "" {
			ConfBase.TimeLocation = ConfBase.Base.TimeLocation
		} else {
			ConfBase.TimeLocation = "Asia/Chongqing"
		}
	}

	logConf := log.Config{
		Level: ConfBase.Log.Level,
		FW: log.ConfigFileWriter{
			On:              ConfBase.Log.FW.On,
			LogPath:         ConfBase.Log.FW.LogPath,
			RotateLogPath:   ConfBase.Log.FW.RotateLogPath,
			RotateWfLogPath: ConfBase.Log.FW.RotateWfLogPath,
			WfLogPath:       ConfBase.Log.FW.WfLogPath,
		},
		CW: log.ConfConsoleWriter{
			Color: ConfBase.Log.CW.Color,
			On:    ConfBase.Log.CW.On,
		},
	}
	if err := log.SetupDefaultLogWithConf(logConf); err != nil {
		panic(err)
	}
	log.SetLayout("2006-01-02 15:04:05.000")
	return nil
}

func InitViperConf() error {
	file, err := os.Open(ConfEnvPath + string(os.PathSeparator))
	if err != nil {
		return err
	}
	fileList, err := file.Readdir(1024)
	if err != nil {
		return err
	}
	for _, f := range fileList {
		if !f.IsDir() {
			conf, err := ioutil.ReadFile(ConfEnvPath + string(os.PathSeparator) + f.Name())
			if err != nil {
				return err
			}
			v := viper.New()
			confType := strings.Split(f.Name(), ".")
			v.SetConfigType(confType[len(confType)-1])
			v.ReadConfig(bytes.NewBuffer(conf))
			if ViperConfMap == nil {
				ViperConfMap = make(map[string]*viper.Viper)
			}
			ViperConfMap[confType[0]] = v
		}
	}
	return nil
}

//获取get配置信息
func GetStringConf(key string) string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return ""
	}
	v, ok := ViperConfMap[keys[0]]
	if !ok {
		return ""
	}
	confString := v.GetString(strings.Join(keys[1:], "."))
	return confString
}

//获取get配置信息
func GetStringMapConf(key string) map[string]interface{} {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetStringMap(strings.Join(keys[1:], "."))
	return conf
}

//获取get配置信息
func GetConf(key string) interface{} {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := ViperConfMap[keys[0]]
	conf := v.Get(strings.Join(keys[1:], "."))
	return conf
}

//获取get配置信息
func GetBoolConf(key string) bool {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return false
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetBool(strings.Join(keys[1:], "."))
	return conf
}

//获取get配置信息
func GetFloat64Conf(key string) float64 {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetFloat64(strings.Join(keys[1:], "."))
	return conf
}

//获取get配置信息
func GetIntConf(key string) int {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetInt(strings.Join(keys[1:], "."))
	return conf
}

//获取get配置信息
func GetStringMapStringConf(key string) map[string]string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetStringMapString(strings.Join(keys[1:len(keys)], "."))
	return conf
}

//获取get配置信息
func GetStringSliceConf(key string) []string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetStringSlice(strings.Join(keys[1:len(keys)], "."))
	return conf
}

//获取get配置信息
func GetTimeConf(key string) time.Time {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return time.Now()
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetTime(strings.Join(keys[1:len(keys)], "."))
	return conf
}

//获取时间阶段长度
func GetDurationConf(key string) time.Duration {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := ViperConfMap[keys[0]]
	conf := v.GetDuration(strings.Join(keys[1:], "."))
	return conf
}

//是否设置了key
func IsSetConf(key string) bool {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return false
	}
	v, ok := ViperConfMap[keys[0]]
	if !ok {
		return false
	}

	conf := v.IsSet(strings.Join(keys[1:], "."))
	return conf
}
