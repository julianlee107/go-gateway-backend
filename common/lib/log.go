package lib

import (
	"fmt"
	"github.com/julianlee107/go-gateway-backend/common/log"
	"strings"
)

const (
	TagUndefined = "_undef"
)

const (
	_tag         = "tag"
	_traceId     = "trace_id"
	_spanId      = "span_id"
	_childSpanId = "child_span_id"
	_tagPrefix   = "_com_"
)

type Logger struct {
}

var Log *Logger

type Trace struct {
	TraceId     string
	SpanId      string
	Caller      string
	SrcMethod   string
	HintCode    int64
	HintContent string
}

type TraceContext struct {
	Trace
	CSpanId string
}

func (l *Logger) Close() {
	log.Close()
}

func (l *Logger) TagInfo(trace *TraceContext, tag string, m map[string]interface{}) {
	m[_tag] = tag
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	log.Info(parseParams(m))
}

func (l *Logger) TagWarning(trace *TraceContext, tag string, m map[string]interface{}) {
	m[_tag] = tag
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	log.Warn(parseParams(m))
}

func (l *Logger) TagDebug(trace *TraceContext, tag string, m map[string]interface{}) {
	m[_tag] = tag
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	log.Debug(parseParams(m))
}

func (l *Logger) TagTrace(trace *TraceContext, tag string, m map[string]interface{}) {
	m[_tag] = tag
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	log.Trace(parseParams(m))
}

func (l *Logger) TagError(trace *TraceContext, tag string, m map[string]interface{}) {
	m[_tag] = tag
	m[_traceId] = trace.TraceId
	m[_childSpanId] = trace.CSpanId
	m[_spanId] = trace.SpanId
	log.Error(parseParams(m))
}

func parseParams(m map[string]interface{}) string {
	var tag = "_undef"
	if _tag, ok := m["tag"]; ok {
		if val, ok := _tag.(string); ok {
			tag = val
		}
	}
	for key, val := range m {
		if key == "tag" {
			continue
		}
		tag = tag + "||" + fmt.Sprintf("%v=%v", key, val)
	}
	tag = strings.Trim(fmt.Sprintf("%q", tag), "\"")
	return tag
}
