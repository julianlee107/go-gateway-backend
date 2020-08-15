package public

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/julianlee107/go-common/lib"
)

//错误日志
func ContextWarning(c context.Context, tag string, m map[string]interface{}) {
	v := c.Value("trace")
	traceContext, ok := v.(*lib.TraceContext)
	if !ok {
		traceContext = lib.NewTrace()
	}
	lib.Log.TagWarning(traceContext, tag, m)
}

//普通日志
func ComLogNotice(c *gin.Context, tag string, m map[string]interface{}) {
	traceContext := GetGinTraceContext(c)
	lib.Log.TagInfo(traceContext, tag, m)
}

// 从gin的Context中获取数据
func GetGinTraceContext(c *gin.Context) *lib.TraceContext {
	// 防御
	if c == nil {
		return lib.NewTrace()
	}
	traceContext, exists := c.Get("trace")
	if exists {
		if tc, ok := traceContext.(*lib.TraceContext); ok {
			return tc
		}
	}
	return lib.NewTrace()
}

// 从Context中获取数据
func GetTraceContext(c context.Context) *lib.TraceContext {
	if c == nil {
		return lib.NewTrace()
	}
	traceContext := c.Value("trace")
	if tc, ok := traceContext.(*lib.TraceContext); ok {
		return tc
	}
	return lib.NewTrace()
}
