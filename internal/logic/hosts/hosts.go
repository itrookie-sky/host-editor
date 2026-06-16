package hosts

import (
	"context"
	"fmt"
	"host-editor/internal/consts"
	"host-editor/internal/model"
	"host-editor/internal/service"
	"strings"

	"github.com/gogf/gf/v2/os/gfile"
)

type sHosts struct {
	dir string
}

func New() *sHosts {
	return &sHosts{}
}

func init() {
	service.RegisterHosts(New())
}

func (s *sHosts) Start(ctx context.Context) (err error) {
	home, err := gfile.Home()
	if err != nil {
		err = fmt.Errorf("get home dir: %w", err)
		return
	}
	dir := gfile.Join(home, consts.HostStoreDirName)
	err = gfile.Mkdir(dir)
	if err != nil {
		err = fmt.Errorf("create store dir: %w", err)
		return
	}
	s.dir = dir

	paths, err := gfile.ScanDirFile(s.dir, "*")
	if err != nil {
		err = fmt.Errorf("read host store dir: %w", err)
		return
	}
	hasHostFile := false
	for _, path := range paths {
		name := gfile.Basename(path)
		if name == consts.DefaultHostFile || strings.HasSuffix(name, consts.HostFileExt) {
			hasHostFile = true
			break
		}
	}
	if hasHostFile {
		return
	}

	data := gfile.GetContents(consts.SystemHostsPath)
	if data == "" && !gfile.Exists(consts.SystemHostsPath) {
		err = fmt.Errorf("read system hosts: file %q does not exist", consts.SystemHostsPath)
		return
	}
	path := gfile.Join(s.dir, consts.DefaultHostFile)
	tmp := path + ".tmp"
	err = gfile.PutContents(tmp, data)
	if err != nil {
		err = fmt.Errorf("write temp file: %w", err)
		return
	}
	err = gfile.Rename(tmp, path)
	if err != nil {
		_ = gfile.RemoveFile(tmp)
		err = fmt.Errorf("rename temp file: %w", err)
		return
	}
	return
}

// ---- public methods implementing IHosts ----

func (s *sHosts) ListHostFiles(ctx context.Context) (res []*model.HostFileInfo, err error) {
	if s.dir == "" {
		err = fmt.Errorf("store directory not initialized")
		return
	}

	paths, err := gfile.ScanDirFile(s.dir, "*")
	if err != nil {
		err = fmt.Errorf("read host store dir: %w", err)
		return
	}

	res = make([]*model.HostFileInfo, 0, len(paths))
	for _, path := range paths {
		name := gfile.Basename(path)
		if name != consts.DefaultHostFile && !strings.HasSuffix(name, consts.HostFileExt) {
			continue
		}
		if name == consts.DefaultHostFile {
			res = append(res, &model.HostFileInfo{Name: consts.DefaultHostFile})
			continue
		}
		res = append(res, &model.HostFileInfo{Name: strings.TrimSuffix(name, consts.HostFileExt)})
	}
	return
}

func (s *sHosts) ReadHostFile(ctx context.Context, name string) (res string, err error) {
	if s.dir == "" {
		err = fmt.Errorf("store directory not initialized")
		return
	}
	if name == "" {
		err = fmt.Errorf("file name cannot be empty")
		return
	}
	if len(name) > consts.MaxHostFileNameLength {
		err = fmt.Errorf("file name too long (max %d characters)", consts.MaxHostFileNameLength)
		return
	}
	if strings.ContainsAny(name, "/\\:\x00") {
		err = fmt.Errorf("file name contains invalid characters")
		return
	}
	if name == "." || name == ".." {
		err = fmt.Errorf("invalid file name")
		return
	}

	path := gfile.Join(s.dir, name+consts.HostFileExt)
	if name == consts.DefaultHostFile {
		path = gfile.Join(s.dir, consts.DefaultHostFile)
	}
	if !gfile.Exists(path) {
		err = fmt.Errorf("read host file: file %q does not exist", name)
		return
	}
	res = gfile.GetContents(path)
	return
}

func (s *sHosts) SaveHostFile(ctx context.Context, req *model.SaveHostFileRequest) (err error) {
	if s.dir == "" {
		err = fmt.Errorf("store directory not initialized")
		return
	}
	if req == nil {
		err = fmt.Errorf("request cannot be nil")
		return
	}
	if req.Name == "" {
		err = fmt.Errorf("file name cannot be empty")
		return
	}
	if len(req.Name) > consts.MaxHostFileNameLength {
		err = fmt.Errorf("file name too long (max %d characters)", consts.MaxHostFileNameLength)
		return
	}
	if strings.ContainsAny(req.Name, "/\\:\x00") {
		err = fmt.Errorf("file name contains invalid characters")
		return
	}
	if req.Name == "." || req.Name == ".." {
		err = fmt.Errorf("invalid file name")
		return
	}

	path := gfile.Join(s.dir, req.Name+consts.HostFileExt)
	if req.Name == consts.DefaultHostFile {
		path = gfile.Join(s.dir, consts.DefaultHostFile)
	}
	tmp := path + ".tmp"
	err = gfile.PutContents(tmp, req.Content)
	if err != nil {
		err = fmt.Errorf("write temp file: %w", err)
		return
	}
	err = gfile.Rename(tmp, path)
	if err != nil {
		_ = gfile.RemoveFile(tmp)
		err = fmt.Errorf("rename temp file: %w", err)
		return
	}
	return
}

func (s *sHosts) CreateHostFile(ctx context.Context, name string) (res *model.HostFileInfo, err error) {
	if s.dir == "" {
		err = fmt.Errorf("store directory not initialized")
		return
	}
	if name == "" {
		err = fmt.Errorf("file name cannot be empty")
		return
	}
	if len(name) > consts.MaxHostFileNameLength {
		err = fmt.Errorf("file name too long (max %d characters)", consts.MaxHostFileNameLength)
		return
	}
	if strings.ContainsAny(name, "/\\:\x00") {
		err = fmt.Errorf("file name contains invalid characters")
		return
	}
	if name == "." || name == ".." {
		err = fmt.Errorf("invalid file name")
		return
	}

	path := gfile.Join(s.dir, name+consts.HostFileExt)
	if name == consts.DefaultHostFile {
		path = gfile.Join(s.dir, consts.DefaultHostFile)
	}
	if gfile.Exists(path) {
		err = fmt.Errorf("file %q already exists", name)
		return
	}
	err = gfile.PutContents(path, "")
	if err != nil {
		err = fmt.Errorf("create host file: %w", err)
		return
	}
	res = &model.HostFileInfo{Name: name}
	return
}

func (s *sHosts) DeleteHostFile(ctx context.Context, name string) (err error) {
	if s.dir == "" {
		err = fmt.Errorf("store directory not initialized")
		return
	}
	if name == "" {
		err = fmt.Errorf("file name cannot be empty")
		return
	}
	if len(name) > consts.MaxHostFileNameLength {
		err = fmt.Errorf("file name too long (max %d characters)", consts.MaxHostFileNameLength)
		return
	}
	if strings.ContainsAny(name, "/\\:\x00") {
		err = fmt.Errorf("file name contains invalid characters")
		return
	}
	if name == "." || name == ".." {
		err = fmt.Errorf("invalid file name")
		return
	}

	path := gfile.Join(s.dir, name+consts.HostFileExt)
	if name == consts.DefaultHostFile {
		path = gfile.Join(s.dir, consts.DefaultHostFile)
	}
	err = gfile.RemoveFile(path)
	if err != nil {
		err = fmt.Errorf("delete host file: %w", err)
		return
	}
	return
}
