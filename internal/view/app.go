package view

import (
	"context"
	"fmt"
	"host-editor/internal/model"
	"host-editor/internal/service"

	"github.com/gogf/gf/v2/os/gctx"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = gctx.WithSpan(ctx, "startup")

	if err := service.Hosts().Init(a.ctx); err != nil {
		fmt.Printf("host-editor: init hosts: %v\n", err)
	}
}

func (a *App) ListHostFiles() ([]model.HostFileInfo, error) {
	return service.Hosts().ListHostFiles()
}

func (a *App) ReadHostFile(name string) (string, error) {
	return service.Hosts().ReadHostFile(name)
}

func (a *App) SaveHostFile(req model.SaveHostFileRequest) error {
	return service.Hosts().SaveHostFile(req)
}

func (a *App) CreateHostFile(name string) (model.HostFileInfo, error) {
	return service.Hosts().CreateHostFile(name)
}

func (a *App) DeleteHostFile(name string) error {
	return service.Hosts().DeleteHostFile(name)
}
