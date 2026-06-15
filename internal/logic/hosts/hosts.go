package hosts

import (
	"context"
	"fmt"
	"host-editor/internal/consts"
	"host-editor/internal/model"
	"host-editor/internal/service"
	"os"
	"path/filepath"
	"strings"
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
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home dir: %w", err)
	}
	dir := filepath.Join(home, consts.HostStoreDirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
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

	data, err := os.ReadFile(consts.SystemHostsPath)
	if err != nil {
		return fmt.Errorf("read system hosts: %w", err)
	}
	return s.saveHostFile(consts.DefaultHostFile, string(data))
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
		return filepath.Join(s.dir, consts.DefaultHostFile)
	}
	return filepath.Join(s.dir, name+consts.HostFileExt)
}

func (s *sHosts) isHostFile(entry os.DirEntry) bool {
	if entry.IsDir() {
		return false
	}
	name := entry.Name()
	return name == consts.DefaultHostFile || strings.HasSuffix(name, consts.HostFileExt)
}

func (s *sHosts) listHostFiles() ([]model.HostFileInfo, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, fmt.Errorf("read host store dir: %w", err)
	}

	var files []model.HostFileInfo
	for _, e := range entries {
		if !s.isHostFile(e) {
			continue
		}
		name := e.Name()
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
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read host file: %w", err)
	}
	return string(data), nil
}

func (s *sHosts) saveHostFile(name, content string) error {
	if err := s.validateHostFileName(name); err != nil {
		return err
	}
	path := s.hostFilePath(name)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}

func (s *sHosts) createHostFile(name string) (model.HostFileInfo, error) {
	if err := s.validateHostFileName(name); err != nil {
		return model.HostFileInfo{}, err
	}
	path := s.hostFilePath(name)
	if _, err := os.Stat(path); err == nil {
		return model.HostFileInfo{}, fmt.Errorf("file %q already exists", name)
	}
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		return model.HostFileInfo{}, fmt.Errorf("create host file: %w", err)
	}
	return model.HostFileInfo{Name: name}, nil
}

func (s *sHosts) deleteHostFile(name string) error {
	if err := s.validateHostFileName(name); err != nil {
		return err
	}
	path := s.hostFilePath(name)
	if err := os.Remove(path); err != nil {
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
