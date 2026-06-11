package logger

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// LogContext 日志功能的上下文
type LogContext struct {
	Log   g.Map
	Start *gtime.Time
	// 文件回溯 日志输出跳过的层级
	skip int
}

// NewLogContext 创建日志上下文
func NewLogContext(optFuncs ...LogContextOptionFunc) *LogContext {
	lctx := &LogContext{
		Log:  g.Map{},
		skip: 1,
	}
	for _, optFunc := range optFuncs {
		optFunc(lctx)
	}
	return lctx
}

// LogContextOptionFunc 日志上下文配置函数
type LogContextOptionFunc func(lctx *LogContext)

// WithCost 日志输出耗时
func WithCost(now *gtime.Time) LogContextOptionFunc {
	return func(lctx *LogContext) {
		lctx.Start = now
	}
}

// formatLog 日志输出格式化
func (lctx *LogContext) formatLog(key string, args ...any) (title string) {
	if len(args) > 0 && len(args)%2 == 0 {
		for i := 0; i < len(args); i += 2 {
			lctx.SetLogValue(args[i].(string), args[i+1])
		}
	}
	costStr := ""
	if lctx.Start != nil {
		cost := gtime.Now().Sub(lctx.Start)
		costStr = fmt.Sprintf("[%v] ", cost.String())
	}
	title = costStr + "%v log: %v"
	return
}

// Skip 文件回溯
func (lctx *LogContext) Skip(line int) *LogContext {
	lctx.skip += line
	return lctx
}

// Debug 日志输出
/*
	输入 args 格式为:
	key1, value1, key2, value2, key3, value3...
*/
func (lctx *LogContext) Debug(ctx context.Context, key string, args ...any) {
	title := lctx.formatLog(key, args...)
	LogAccess().Skip(lctx.skip).Debugf(ctx, title, key, String(lctx.Log))
}

// Info info级日志
func (lctx *LogContext) Info(ctx context.Context, key string, args ...any) {
	title := lctx.formatLog(key, args...)
	LogAccess().Skip(lctx.skip).Infof(ctx, title, key, String(lctx.Log))
}

// SetLogValue 更新日志 单个字段
func (lctx *LogContext) SetLogValue(key string, value any) {
	if lctx.Log == nil {
		return
	}
	lctx.Log[key] = value
}

// SetLogErr 更新日志 错误信息
func (lctx *LogContext) SetLogErr(e *error) {
	if e != nil && *e != nil {
		err := *e
		rCode := gerror.Code(err)
		if rCode == gcode.CodeNil {
			lctx.SetLogValue("err", err.Error())
			return
		}
		lctx.UpdateLogMap(g.Map{
			"err": g.Map{
				"code":   rCode.Code(),
				"msg":    err.Error(),
				"detail": rCode.Detail(),
			},
		})
	}
}

// UpdateLogMap 更新日志
func (lctx *LogContext) UpdateLogMap(m g.Map) {
	if lctx.Log == nil || len(m) == 0 {
		return
	}
	for k, v := range m {
		lctx.Log[k] = v
	}
}
