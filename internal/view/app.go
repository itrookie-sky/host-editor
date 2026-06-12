package view

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/os/gctx"
)

type App struct {
	ctx  context.Context
	dir  string
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

	files, _ := listHostFiles(dir)
	if len(files) == 0 {
		if _, err := createHostFile(dir, defaultFile); err != nil {
			fmt.Printf("host-editor: create default hosts: %v\n", err)
		}
	}
}

func (a *App) ListHostFiles() ([]HostFileInfo, error) {
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

func (a *App) SaveHostFile(req SaveHostFileRequest) error {
	if a.dir == "" {
		return fmt.Errorf("store directory not initialized")
	}
	return saveHostFile(a.dir, req.Name, req.Content)
}

func (a *App) CreateHostFile(name string) (HostFileInfo, error) {
	if a.dir == "" {
		return HostFileInfo{}, fmt.Errorf("store directory not initialized")
	}
	return createHostFile(a.dir, name)
}

func (a *App) DeleteHostFile(name string) error {
	if a.dir == "" {
		return fmt.Errorf("store directory not initialized")
	}
	return deleteHostFile(a.dir, name)
}
