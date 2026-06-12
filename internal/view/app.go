package view

import (
	"context"
	"fmt"
	"host-editor/internal/consts"
	"host-editor/internal/model"

	"github.com/gogf/gf/v2/os/gctx"
)

type App struct {
	ctx context.Context
	dir string
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = gctx.WithSpan(ctx, "startup")

	dir, err := hostStoreDir()
	if err != nil {
		fmt.Printf("host-editor: resolve store dir: %v\n", err)
		return
	}
	if err := ensureStoreDir(dir); err != nil {
		fmt.Printf("host-editor: create store dir: %v\n", err)
		return
	}
	a.dir = dir

	err = ensureDefaultHostFile(dir, consts.SystemHostsPath)
	if err != nil {
		fmt.Printf("host-editor: create default hosts: %v\n", err)
	}
}

func (a *App) ListHostFiles() ([]model.HostFileInfo, error) {
	if a.dir == "" {
		return nil, fmt.Errorf("store directory not initialized")
	}
	return listHostFiles(a.dir)
}

func (a *App) ReadHostFile(name string) (string, error) {
	if a.dir == "" {
		return "", fmt.Errorf("store directory not initialized")
	}
	return readHostFile(a.dir, name)
}

func (a *App) SaveHostFile(req model.SaveHostFileRequest) error {
	if a.dir == "" {
		return fmt.Errorf("store directory not initialized")
	}
	return saveHostFile(a.dir, req.Name, req.Content)
}

func (a *App) CreateHostFile(name string) (model.HostFileInfo, error) {
	if a.dir == "" {
		return model.HostFileInfo{}, fmt.Errorf("store directory not initialized")
	}
	return createHostFile(a.dir, name)
}

func (a *App) DeleteHostFile(name string) error {
	if a.dir == "" {
		return fmt.Errorf("store directory not initialized")
	}
	return deleteHostFile(a.dir, name)
}
