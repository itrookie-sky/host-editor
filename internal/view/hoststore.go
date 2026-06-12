package view

import (
	"fmt"
	"host-editor/internal/consts"
	"host-editor/internal/model"
	"os"
	"path/filepath"
	"strings"
)

func hostStoreDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	return filepath.Join(home, consts.HostStoreDirName), nil
}

func ensureStoreDir(dir string) error {
	return os.MkdirAll(dir, 0o755)
}

func validateHostFileName(name string) error {
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

func hostFilePath(dir, name string) string {
	if name == consts.DefaultHostFile {
		return filepath.Join(dir, consts.DefaultHostFile)
	}
	return filepath.Join(dir, name+consts.HostFileExt)
}

func isHostFile(entry os.DirEntry) bool {
	if entry.IsDir() {
		return false
	}
	name := entry.Name()
	return name == consts.DefaultHostFile || strings.HasSuffix(name, consts.HostFileExt)
}

func listHostFiles(dir string) ([]model.HostFileInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read host store dir: %w", err)
	}

	var files []model.HostFileInfo
	for _, e := range entries {
		if !isHostFile(e) {
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

func readHostFile(dir, name string) (string, error) {
	if err := validateHostFileName(name); err != nil {
		return "", err
	}
	path := hostFilePath(dir, name)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read host file: %w", err)
	}
	return string(data), nil
}

func saveHostFile(dir, name, content string) error {
	if err := validateHostFileName(name); err != nil {
		return err
	}
	path := hostFilePath(dir, name)
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

func createHostFile(dir, name string) (model.HostFileInfo, error) {
	if err := validateHostFileName(name); err != nil {
		return model.HostFileInfo{}, err
	}
	path := hostFilePath(dir, name)
	if _, err := os.Stat(path); err == nil {
		return model.HostFileInfo{}, fmt.Errorf("file %q already exists", name)
	}
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		return model.HostFileInfo{}, fmt.Errorf("create host file: %w", err)
	}
	return model.HostFileInfo{Name: name}, nil
}

func ensureDefaultHostFile(dir, sourcePath string) error {
	files, err := listHostFiles(dir)
	if err != nil {
		return err
	}
	if len(files) > 0 {
		return nil
	}

	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("read system hosts: %w", err)
	}
	return saveHostFile(dir, consts.DefaultHostFile, string(data))
}

func deleteHostFile(dir, name string) error {
	if err := validateHostFileName(name); err != nil {
		return err
	}
	path := hostFilePath(dir, name)
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("delete host file: %w", err)
	}
	return nil
}
