package view

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	dirName      = ".host-editor"
	defaultFile  = "hosts"
	fileExt      = ".hosts"
	maxFileNameLen = 64
)

type HostFileInfo struct {
	Name    string `json:"name"`
	IsDirty bool   `json:"isDirty"`
}

type SaveHostFileRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func hostStoreDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	return filepath.Join(home, dirName), nil
}

func ensureStoreDir(dir string) error {
	return os.MkdirAll(dir, 0o755)
}

func validateHostFileName(name string) error {
	if name == "" {
		return fmt.Errorf("file name cannot be empty")
	}
	if len(name) > maxFileNameLen {
		return fmt.Errorf("file name too long (max %d characters)", maxFileNameLen)
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
	if name == defaultFile {
		return filepath.Join(dir, defaultFile)
	}
	return filepath.Join(dir, name+fileExt)
}

func isHostFile(entry os.DirEntry) bool {
	if entry.IsDir() {
		return false
	}
	name := entry.Name()
	return name == defaultFile || strings.HasSuffix(name, fileExt)
}

func listHostFiles(dir string) ([]HostFileInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read host store dir: %w", err)
	}

	var files []HostFileInfo
	for _, e := range entries {
		if !isHostFile(e) {
			continue
		}
		name := e.Name()
		if name == defaultFile {
			files = append(files, HostFileInfo{Name: defaultFile})
		} else {
			files = append(files, HostFileInfo{Name: strings.TrimSuffix(name, fileExt)})
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

func createHostFile(dir, name string) (HostFileInfo, error) {
	if err := validateHostFileName(name); err != nil {
		return HostFileInfo{}, err
	}
	path := hostFilePath(dir, name)
	if _, err := os.Stat(path); err == nil {
		return HostFileInfo{}, fmt.Errorf("file %q already exists", name)
	}
	if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
		return HostFileInfo{}, fmt.Errorf("create host file: %w", err)
	}
	return HostFileInfo{Name: name}, nil
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
