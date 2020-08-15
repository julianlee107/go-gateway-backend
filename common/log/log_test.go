package log

import (
	"testing"
	"time"
)

//测试日志实例打点
func TestLogInstance(t *testing.T) {
	nlog := NewLogger()
	logConf := Config{
		Level: "trace",
		FW: ConfigFileWriter{
			On:              true,
			LogPath:         "./%Y-%M-%D_log_test.log",
			RotateLogPath:   "./%Y-%M-%D_log_test.log",
			WfLogPath:       "./%Y-%M-%D_log_test.wf.log",
			RotateWfLogPath: "./%Y-%M-%D_log_test.wf.log",
		},
		CW: ConfConsoleWriter{
			On:    true,
			Color: true,
		},
	}
	SetupInstanceWithConf(logConf, nlog)
	nlog.Warn("test message")
	nlog.Close()
	time.Sleep(10 * time.Second)
}
