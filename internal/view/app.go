package view

import (
	"context"
	"host-editor/internal/service"
	"host-editor/utility/logger"
	"sync"

	"github.com/gogf/gf/v2/os/gctx"
)

var (
	app  *App
	once sync.Once
)

type App struct {
	ctx   context.Context
	hosts *Hosts
}

// NewApp 创建 App 实例（单例模式）
func NewApp() *App {
	once.Do(func() {
		app = &App{
			hosts: &Hosts{},
		}
	})
	return app
}

// GetBind 获取绑定实例
func GetBind() []interface{} {
	app := NewApp()
	return []interface{}{
		app,
		app.hosts,
	}
}

// Startup 应用启动
func (a *App) Startup(ctx context.Context) {
	ctx = gctx.WithSpan(ctx, "startup")
	a.ctx = ctx
	a.hosts.SetCtx(ctx)
	// 前置服务
	err := service.Hosts().Start(ctx)
	if err != nil {
		logger.Errorf(ctx, "service.hosts.start: %v", err)
		return
	}
	logger.Infof(ctx, "app.startup")
}

// Shutdown 退出
func (a *App) Shutdown(ctx context.Context) {
	logger.Info(a.ctx, "app.shutdown")
}
