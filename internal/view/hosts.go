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
func (h *Hosts) SetCtx(ctx context.Context) {
	h.ctx = ctx
}

func (h *Hosts) ListHostFiles() (res []*model.HostFileInfo, err error) {
	return service.Hosts().ListHostFiles(h.ctx)
}

func (h *Hosts) ReadHostFile(name string) (res string, err error) {
	return service.Hosts().ReadHostFile(h.ctx, name)
}

func (h *Hosts) SaveHostFile(req model.SaveHostFileRequest) (err error) {
	return service.Hosts().SaveHostFile(h.ctx, &req)
}

func (h *Hosts) CreateHostFile(name string) (res *model.HostFileInfo, err error) {
	return service.Hosts().CreateHostFile(h.ctx, name)
}

func (h *Hosts) DeleteHostFile(name string) (err error) {
	return service.Hosts().DeleteHostFile(h.ctx, name)
}
