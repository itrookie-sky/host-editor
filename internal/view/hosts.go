package view

import (
	"context"
	"host-editor/internal/model"
	"host-editor/internal/service"
)

type Hosts struct {
	ctx context.Context
}

// SetCtx 设置上下文
func (a *Hosts) SetCtx(ctx context.Context) {
	a.ctx = ctx
}

func (a *Hosts) ListHostFiles() ([]model.HostFileInfo, error) {
	return service.Hosts().ListHostFiles()
}

func (a *Hosts) ReadHostFile(name string) (string, error) {
	return service.Hosts().ReadHostFile(name)
}

func (a *Hosts) SaveHostFile(req model.SaveHostFileRequest) error {
	return service.Hosts().SaveHostFile(req)
}

func (a *Hosts) CreateHostFile(name string) (model.HostFileInfo, error) {
	return service.Hosts().CreateHostFile(name)
}

func (a *Hosts) DeleteHostFile(name string) error {
	return service.Hosts().DeleteHostFile(name)
}
