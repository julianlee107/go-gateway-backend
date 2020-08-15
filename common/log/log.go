package log

import (
	"fmt"
	"log"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var (
	LEVEL_FLAGS = [...]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

const (
	TRACE = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

const tunnel_size_default = 1024

type Record struct {
	time  string
	code  string
	info  string
	level int
}

func (r *Record) String() string {
	return fmt.Sprintf("[%s][%s][%s] %s\n", LEVEL_FLAGS[r.level], r.time, r.code, r.info)
}

type Writer interface {
	Init() error
	Write(*Record) error
}

type Rotater interface {
	Rotate() error
	SetPathPattern(string) error
}

type Flusher interface {
	Flush() error
}

type Logger struct {
	writers     []Writer
	tunnel      chan *Record
	level       int
	lastTime    int64
	lastTimeStr string
	c           chan bool
	layout      string
	recordPool  *sync.Pool
}

var (
	defaultLogger *Logger
	sw            = false //默认启动标识
)

func NewLogger() *Logger {
	if defaultLogger != nil && sw == false {
		sw = true
		return defaultLogger
	}
	logger := new(Logger)
	logger.writers = []Writer{}
	logger.tunnel = make(chan *Record, tunnel_size_default)
	logger.c = make(chan bool, 2)
	logger.level = DEBUG
	logger.layout = "2006/01/02 15:04:05"
	logger.recordPool = &sync.Pool{New: func() interface{} {
		return &Record{}
	}}

	go bootstrapLogWriter(logger)

	return logger
}

func (l *Logger) Register(w Writer) {
	if err := w.Init(); err != nil {
		panic(err)
	}
	l.writers = append(l.writers, w)
}

func (l *Logger) SetLevel(level int) {
	l.level = level
}

func (l *Logger) SetLayout(layout string) {
	l.layout = layout
}
func (l *Logger) Trace(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(TRACE, fmt, args...)
}

func (l *Logger) Debug(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(DEBUG, fmt, args...)
}

func (l *Logger) Warn(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(WARNING, fmt, args...)
}

func (l *Logger) Info(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(INFO, fmt, args...)
}

func (l *Logger) Error(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(ERROR, fmt, args...)
}

func (l *Logger) Fatal(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(FATAL, fmt, args...)
}

func (l *Logger) Close() {
	close(l.tunnel)
	<-l.c
	for _, w := range l.writers {
		if f, ok := w.(Flusher); ok {
			if err := f.Flush(); err != nil {
				log.Println(err)
			}
		}
	}
}

func (l *Logger) deliverRecordToWriter(level int, format string, args ...interface{}) {
	var inf, code string
	if level < l.level {
		return
	}

	if format != "" {
		inf = fmt.Sprintf(format, args...)
	} else {
		inf = fmt.Sprint(args...)
	}
	// source code, file and line num
	_, file, line, ok := runtime.Caller(2)
	if ok {
		code = path.Base(file) + ":" + strconv.Itoa(line)
	}
	//	format time
	now := time.Now()
	if now.Unix() != l.lastTime {
		l.lastTime = now.Unix()
		l.lastTimeStr = now.Format(l.layout)
	}
	record := l.recordPool.Get().(*Record)
	record.info = inf
	record.code = code
	record.time = l.lastTimeStr
	record.level = level

	l.tunnel <- record
}

func bootstrapLogWriter(logger *Logger) {
	if logger == nil {
		panic("logger is nil")
	}
	var (
		record *Record
		ok     bool
	)
	if record, ok = <-logger.tunnel; !ok {
		logger.c <- true
		return
	}
	for _, w := range logger.writers {
		if err := w.Write(record); err != nil {
			log.Println(err)
		}
	}

	flushTimer := time.NewTimer(time.Millisecond * 500)
	rotateTimer := time.NewTimer(time.Second * 10)

	for {
		select {
		case record, ok := <-logger.tunnel:
			if !ok {
				logger.c <- true
				return
			}
			for _, w := range logger.writers {
				if err := w.Write(record); err != nil {
					log.Println(err)
				}
			}
			logger.recordPool.Put(record)
		case <-flushTimer.C:
			for _, w := range logger.writers {
				if f, ok := w.(Flusher); ok {
					if err := f.Flush(); err != nil {
						log.Println(err)
					}
				}
			}
			flushTimer.Reset(time.Millisecond * 500)
		case <-rotateTimer.C:
			for _, w := range logger.writers {
				if f, ok := w.(Rotater); ok {
					if err := f.Rotate(); err != nil {
						log.Println(err)
					}
				}
			}
			rotateTimer.Reset(time.Second * 10)

		}
	}
}

// outside
func SetLevel(lvl int) {
	defaultLoggerInit()
	defaultLogger.level = lvl
}

func SetLayout(layout string) {
	defaultLoggerInit()
	defaultLogger.layout = layout
}

func Trace(fmt string, args ...interface{}) {
	defaultLoggerInit()
	defaultLogger.deliverRecordToWriter(TRACE, fmt, args...)
}

func Debug(fmt string, args ...interface{}) {
	defaultLoggerInit()
	defaultLogger.deliverRecordToWriter(DEBUG, fmt, args...)
}

func Warn(fmt string, args ...interface{}) {
	defaultLoggerInit()
	defaultLogger.deliverRecordToWriter(WARNING, fmt, args...)
}

func Info(fmt string, args ...interface{}) {
	defaultLoggerInit()
	defaultLogger.deliverRecordToWriter(INFO, fmt, args...)
}

func Error(fmt string, args ...interface{}) {
	defaultLoggerInit()
	defaultLogger.deliverRecordToWriter(ERROR, fmt, args...)
}

func Fatal(fmt string, args ...interface{}) {
	defaultLoggerInit()
	defaultLogger.deliverRecordToWriter(FATAL, fmt, args...)
}

func Register(w Writer) {
	defaultLoggerInit()
	defaultLogger.Register(w)
}

func Close() {
	defaultLoggerInit()
	defaultLogger.Close()
	defaultLogger = nil
	sw = false
}

func defaultLoggerInit() {
	if sw == false {
		defaultLogger = NewLogger()
	}
}

func SetupDefaultLogWithConf(c Config) (err error) {
	defaultLoggerInit()
	return SetupInstanceWithConf(c, defaultLogger)
}
