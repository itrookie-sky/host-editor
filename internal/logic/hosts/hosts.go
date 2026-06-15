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

func (s *sHosts) Start() {}

// Init initializes the host file store directory and ensures the default host file
// exists (by copying from the system hosts file if the store is empty).
func (s *sHosts) Init(_ context.Context) error {
	home, err := gfile.Home()
	if err != nil {
		return fmt.Errorf("get home dir: %w", err)
	}
	dir := gfile.Join(home, consts.HostStoreDirName)
	if err := gfile.Mkdir(dir); err != nil {
		return fmt.Errorf("create store dir: %w", err)
	}
	s.dir = dir

	files, err := s.listHostFiles()
	if err != nil {
		return err
	}
	if len(files) > 0 {
		return nil
	}

	data := gfile.GetContents(consts.SystemHostsPath)
	if data == "" && !gfile.Exists(consts.SystemHostsPath) {
		return fmt.Errorf("read system hosts: file %q does not exist", consts.SystemHostsPath)
	}
	return s.saveHostFile(consts.DefaultHostFile, data)
}

func (s *sHosts) validateHostFileName(name string) error {
	if name == "" {
		return fmt.Errorf("file name cannot be empty")
	}
	if len(name) > consts.MaxHostFileNameLength {
		return fmt.Errorf("file name too long (max %d characters)", consts.MaxHostFileNameLength)
	}
	if strings.ContainsAny(name, "/\\:\x00") {
		return fmt.Errorf("file name contains invalid characters")
	}
	if name == "." || name == ".." {
		return fmt.Errorf("invalid file name")
	}
	return nil
}

func (s *sHosts) hostFilePath(name string) string {
	if name == consts.DefaultHostFile {
		return gfile.Join(s.dir, consts.DefaultHostFile)
	}
	return gfile.Join(s.dir, name+consts.HostFileExt)
}

func (s *sHosts) isHostFile(path string) bool {
	name := gfile.Basename(path)
	return name == consts.DefaultHostFile || strings.HasSuffix(name, consts.HostFileExt)
}

func (s *sHosts) listHostFiles() ([]model.HostFileInfo, error) {
	paths, err := gfile.ScanDirFile(s.dir, "*")
	if err != nil {
		return nil, fmt.Errorf("read host store dir: %w", err)
	}

	var files []model.HostFileInfo
	for _, path := range paths {
		if !s.isHostFile(path) {
			continue
		}
		name := gfile.Basename(path)
		if name == consts.DefaultHostFile {
			files = append(files, model.HostFileInfo{Name: consts.DefaultHostFile})
		} else {
			files = append(files, model.HostFileInfo{Name: strings.TrimSuffix(name, consts.HostFileExt)})
		}
	}
	return files, nil
}

func (s *sHosts) readHostFile(name string) (string, error) {
	if err := s.validateHostFileName(name); err != nil {
		return "", err
	}
	path := s.hostFilePath(name)
	if !gfile.Exists(path) {
		return "", fmt.Errorf("read host file: file %q does not exist", name)
	}
	return gfile.GetContents(path), nil
}

func (s *sHosts) saveHostFile(name, content string) error {
	if err := s.validateHostFileName(name); err != nil {
		return err
	}
	path := s.hostFilePath(name)
	tmp := path + ".tmp"
	if err := gfile.PutContents(tmp, content); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := gfile.Rename(tmp, path); err != nil {
		_ = gfile.RemoveFile(tmp)
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}

func (s *sHosts) createHostFile(name string) (model.HostFileInfo, error) {
	if err := s.validateHostFileName(name); err != nil {
		return model.HostFileInfo{}, err
	}
	path := s.hostFilePath(name)
	if gfile.Exists(path) {
		return model.HostFileInfo{}, fmt.Errorf("file %q already exists", name)
	}
	if err := gfile.PutContents(path, ""); err != nil {
		return model.HostFileInfo{}, fmt.Errorf("create host file: %w", err)
	}
	return model.HostFileInfo{Name: name}, nil
}

func (s *sHosts) deleteHostFile(name string) error {
	if err := s.validateHostFileName(name); err != nil {
		return err
	}
	path := s.hostFilePath(name)
	if err := gfile.RemoveFile(path); err != nil {
		return fmt.Errorf("delete host file: %w", err)
	}
	return nil
}

// ---- public methods implementing IHosts ----

func (s *sHosts) ListHostFiles() ([]model.HostFileInfo, error) {
	if s.dir == "" {
		return nil, fmt.Errorf("store directory not initialized")
	}
	return s.listHostFiles()
}

func (s *sHosts) ReadHostFile(name string) (string, error) {
	if s.dir == "" {
		return "", fmt.Errorf("store directory not initialized")
	}
	return s.readHostFile(name)
}

func (s *sHosts) SaveHostFile(req model.SaveHostFileRequest) error {
	if s.dir == "" {
		return fmt.Errorf("store directory not initialized")
	}
	return s.saveHostFile(req.Name, req.Content)
}

func (s *sHosts) CreateHostFile(name string) (model.HostFileInfo, error) {
	if s.dir == "" {
		return model.HostFileInfo{}, fmt.Errorf("store directory not initialized")
	}
	return s.createHostFile(name)
}

func (s *sHosts) DeleteHostFile(name string) error {
	if s.dir == "" {
		return fmt.Errorf("store directory not initialized")
	}
	return s.deleteHostFile(name)
}
