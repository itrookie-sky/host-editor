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
		Start()
		Init(ctx context.Context) error
		ListHostFiles() ([]model.HostFileInfo, error)
		ReadHostFile(name string) (string, error)
		SaveHostFile(req model.SaveHostFileRequest) error
		CreateHostFile(name string) (model.HostFileInfo, error)
		DeleteHostFile(name string) error
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
