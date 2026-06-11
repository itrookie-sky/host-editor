package logger

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
)

// Log 获取log对象
func Log(name ...string) *glog.Logger {
	if len(name) > 0 && name[0] != "" {
		return g.Log(name[0]).Skip(1).Line()
	}
	return g.Log().Skip(1).Line()
}

func LogAccess() *glog.Logger {
	return Log("access")
}

func LogErr() *glog.Logger {
	return Log("error")
}

func Print(ctx context.Context, v ...any) {
	Log("access").Print(ctx, v...)
}

func Printf(ctx context.Context, format string, v ...any) {
	Log("access").Printf(ctx, format, v...)
}

func Info(ctx context.Context, v ...any) {
	Log("access").Info(ctx, v...)
}

func Infof(ctx context.Context, format string, v ...any) {
	Log("access").Infof(ctx, format, v...)
}

func Debug(ctx context.Context, v ...any) {
	Log("access").Debug(ctx, v...)
}

func Debugf(ctx context.Context, format string, v ...any) {
	Log("access").Debugf(ctx, format, v...)
}

func Cost(ctx context.Context, start *gtime.Time, v ...any) {
	end := gtime.Now()
	diff := end.Sub(start)
	v1 := []any{
		fmt.Sprintf("[ 耗时%.3fs ] ", diff.Seconds()),
	}
	if len(v) > 0 {
		v1 = append(v1, v...)
	}
	Log("access").Debug(ctx, v1...)
}

func Costf(ctx context.Context, start *gtime.Time, format string, v ...any) {
	end := gtime.Now()
	diff := end.Sub(start)
	Log("access").Debugf(ctx, fmt.Sprintf("[ 耗时%.3fs ] ", diff.Seconds())+format, v...)
}

func Notice(ctx context.Context, v ...any) {
	Log("access").Notice(ctx, v...)
}

func Noticef(ctx context.Context, format string, v ...any) {
	Log("access").Noticef(ctx, format, v...)
}

func Warning(ctx context.Context, v ...any) {
	Log("error").Warning(ctx, v...)
}

func Warningf(ctx context.Context, format string, v ...any) {
	Log("error").Warningf(ctx, format, v...)
}

func Error(ctx context.Context, v ...any) {
	Log("error").Error(ctx, v...)
}

func Errorf(ctx context.Context, format string, v ...any) {
	Log("error").Errorf(ctx, format, v...)
}

func Fatal(ctx context.Context, v ...any) {
	Log("error").Fatal(ctx, v...)
}

func Fatalf(ctx context.Context, format string, v ...any) {
	Log("error").Fatalf(ctx, format, v...)
}

func Critical(ctx context.Context, v ...any) {
	Log("error").Critical(ctx, v...)
}

func Criticalf(ctx context.Context, format string, v ...any) {
	Log("error").Criticalf(ctx, format, v...)
}

func Panic(ctx context.Context, v ...any) {
	Log("error").Panic(ctx, v...)
}

func Panicf(ctx context.Context, format string, v ...any) {
	Log("error").Panicf(ctx, format, v...)
}
