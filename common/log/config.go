package log

import (
	"errors"
)

type ConfigFileWriter struct {
	On              bool   `toml:"On"`
	LogPath         string `toml:"LogPath"`
	RotateLogPath   string `toml:"RotateLogPath"`
	WfLogPath       string `toml:"WfLogPath"`
	RotateWfLogPath string `toml:"RotateWfLogPath"`
}

type ConfConsoleWriter struct {
	On    bool `toml:"On"`
	Color bool `toml:"On"`
}

type Config struct {
	Level string            `toml:"LogLevel"`
	FW    ConfigFileWriter  `toml:"FileWriter"`
	CW    ConfConsoleWriter `toml:"ConsoleWriter"`
}

func SetupInstanceWithConf(c Config, logger *Logger) (err error) {
	if c.FW.On {
		if len(c.FW.LogPath) > 0 {
			w := NewFileWriter()
			w.SetFileName(c.FW.LogPath)
			w.SetPathPattern(c.FW.RotateLogPath)
			w.SetLogLevelCeil(TRACE)
			if len(c.FW.WfLogPath) > 0 {
				w.SetLogLevelCeil(INFO)
			} else {
				w.SetLogLevelCeil(ERROR)
			}
			logger.Register(w)
		}
		if len(c.FW.WfLogPath) > 0 {
			wfw := NewFileWriter()
			wfw.SetFileName(c.FW.WfLogPath)
			wfw.SetPathPattern(c.FW.RotateWfLogPath)
			wfw.SetLogLevelFloor(WARNING)
			wfw.SetLogLevelCeil(ERROR)
			logger.Register(wfw)
		}
	}
	if c.CW.On {
		//w := NewConsoleWriter()
		//w.SetColor(lc.CW.Color)
		//logger.Register(w)
	}
	switch c.Level {
	case "trace":
		logger.SetLevel(TRACE)

	case "debug":
		logger.SetLevel(DEBUG)

	case "info":
		logger.SetLevel(INFO)

	case "warning":
		logger.SetLevel(WARNING)

	case "error":
		logger.SetLevel(ERROR)

	case "fatal":
		logger.SetLevel(FATAL)

	default:
		err = errors.New("invalid log level")
	}
	return
}
