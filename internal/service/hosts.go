// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"host-editor/internal/model"
)

type (
	IHosts interface {
		Start(ctx context.Context) (err error)
		ListHostFiles(ctx context.Context) (res []*model.HostFileInfo, err error)
		ReadHostFile(ctx context.Context, name string) (res string, err error)
		SaveHostFile(ctx context.Context, req *model.SaveHostFileRequest) (err error)
		CreateHostFile(ctx context.Context, name string) (res *model.HostFileInfo, err error)
		DeleteHostFile(ctx context.Context, name string) (err error)
	}
)

var (
	localHosts IHosts
)

func Hosts() IHosts {
	if localHosts == nil {
		panic("implement not found for interface IHosts, forgot register?")
	}
	return localHosts
}

func RegisterHosts(i IHosts) {
	localHosts = i
}
